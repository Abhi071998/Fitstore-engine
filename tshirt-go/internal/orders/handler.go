package orders

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// undefinedTable reports whether err is Postgres error 42P01 (relation does
// not exist) — expected while fitstore-core hasn't deployed its order tables yet.
func undefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}

type Handler struct {
	DB *gorm.DB
}

// CustomerPendingOrders bundles every pending-approval order placed by one customer.
type CustomerPendingOrders struct {
	CustUserID uint64  `json:"cust_user_id"`
	Orders     []Order `json:"orders"`
}

// GET /api/orders/pending (Protected)
// Fetches every order awaiting admin approval and groups them by customer.
func (h *Handler) GetPendingOrders(c echo.Context) error {
	var pending []Order
	if err := h.DB.
		Preload("Items.ProductSize.Product.Category").
		Where("status = ?", "pending_approval").
		Order("cust_user_id, created_at").
		Find(&pending).Error; err != nil {
		if undefinedTable(err) {
			return c.JSON(http.StatusOK, []CustomerPendingOrders{})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch pending orders"})
	}

	grouped := make(map[uint64][]Order)
	var custIDs []uint64
	for _, o := range pending {
		if _, seen := grouped[o.CustUserID]; !seen {
			custIDs = append(custIDs, o.CustUserID)
		}
		grouped[o.CustUserID] = append(grouped[o.CustUserID], o)
	}
	sort.Slice(custIDs, func(i, j int) bool { return custIDs[i] < custIDs[j] })

	result := make([]CustomerPendingOrders, 0, len(custIDs))
	for _, id := range custIDs {
		result = append(result, CustomerPendingOrders{CustUserID: id, Orders: grouped[id]})
	}

	return c.JSON(http.StatusOK, result)
}

// RejectOrderDTO carries the admin's reason when declining an order.
type RejectOrderDTO struct {
	Comment string `json:"comment"`
}

// BulkApproveDTO carries the set of orders to approve together.
type BulkApproveDTO struct {
	OrderIDs []uint64 `json:"order_ids"`
}

// decideOrder is the shared guard for approve/reject: only a
// still-pending order can be decided, and decided_at is stamped once.
func (h *Handler) decideOrder(orderID uint64, newStatus string, comment *string) (*Order, error) {
	result := h.DB.Model(&Order{}).
		Where("id = ? AND status = ?", orderID, "pending_approval").
		Updates(map[string]interface{}{
			"status":        newStatus,
			"decided_at":    time.Now(),
			"admin_comment": comment,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var order Order
	if err := h.DB.Preload("Items.ProductSize.Product.Category").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// PUT /api/orders/:id/approve (Protected)
func (h *Handler) ApproveOrder(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID format"})
	}

	order, err := h.decideOrder(orderID, "approved", nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found or already decided"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to approve order"})
	}

	return c.JSON(http.StatusOK, order)
}

// PUT /api/orders/:id/reject (Protected)
func (h *Handler) RejectOrder(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order ID format"})
	}

	dto := new(RejectOrderDTO)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}
	if dto.Comment == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "A comment is required to reject an order"})
	}

	order, err := h.decideOrder(orderID, "rejected", &dto.Comment)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found or already decided"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reject order"})
	}

	return c.JSON(http.StatusOK, order)
}

// POST /api/orders/bulk-approve (Protected)
// Approves every listed order that is still pending_approval in one shot.
// Any IDs that are missing, already decided, or don't belong to a pending
// order are silently skipped and reported back under "skipped".
func (h *Handler) BulkApproveOrders(c echo.Context) error {
	dto := new(BulkApproveDTO)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}
	if len(dto.OrderIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "order_ids must not be empty"})
	}

	var approvedIDs []uint64
	tx := h.DB.Begin()
	for _, id := range dto.OrderIDs {
		result := tx.Model(&Order{}).
			Where("id = ? AND status = ?", id, "pending_approval").
			Updates(map[string]interface{}{
				"status":     "approved",
				"decided_at": time.Now(),
			})
		if result.Error != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bulk approve orders"})
		}
		if result.RowsAffected > 0 {
			approvedIDs = append(approvedIDs, id)
		}
	}
	tx.Commit()

	skipped := make([]uint64, 0)
	approvedSet := make(map[uint64]bool, len(approvedIDs))
	for _, id := range approvedIDs {
		approvedSet[id] = true
	}
	for _, id := range dto.OrderIDs {
		if !approvedSet[id] {
			skipped = append(skipped, id)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"approved": approvedIDs,
		"skipped":  skipped,
	})
}
