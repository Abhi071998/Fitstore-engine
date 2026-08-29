# Database Schema Reference

> ⚠️ **Schema ownership has moved to Prisma.** As of the removal of `AutoMigrate`
> from `cmd/api/main.go`, this Go service no longer creates or alters any table.
> **All 8 tables below — not just `orders`/`order_items` — are now owned and
> migrated exclusively by fitstore-core's Prisma schema.** This doc is the
> handoff spec: whoever maintains `schema.prisma` should match these columns/
> types exactly (the existing tables were originally created by GORM, so
> `prisma db pull` against the live DB first is the safest way to adopt them
> without an unintended drop/alter).

This is a **read-only reference doc**, not a source of truth. The real model/DTO
definitions live in their own files (linked per model below) — if you change a
struct here, update this file to match, but the struct/Prisma schema wins if
they ever disagree.

Covers every GORM model in this repo, including the two tables that were
already Prisma-owned (`orders`, `order_items`) before this change.

---

## 1. `User`
**Defined in:** `internal/models/base.go`
**Table:** `users` (default pluralization, no `TableName()` override)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| Email | `string` | `unique;not null;type:varchar(255)` |
| Name | `string` | `not null;type:varchar(255)` |
| Password | `string` | `not null;type:varchar(255)` (hidden from JSON via `json:"-"`) |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |
| DeletedAt | `gorm.DeletedAt` | `index` (soft delete) |

**Relations:** none.

**DTOs:** `RegisterDTO { Email, Name, Password string }`, `LoginDTO { Email, Password string }` — same file.

---

## 2. `Category`
**Defined in:** `internal/models/base.go`
**Table:** `categories` (default pluralization)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| Name | `string` | `unique;not null;type:varchar(255)` |
| ImageURL | `string` | *(none)* |
| CategoryTypeID | `*uint` | `index` — nullable FK → `category_types.id` |
| CategoryType | `*CategoryType` | `foreignKey:CategoryTypeID` (belongs-to) |
| Products | `[]Product` | `foreignKey:CategoryID` (has-many) |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |
| DeletedAt | `gorm.DeletedAt` | `index` (soft delete) |

**Relations:**
- Belongs to `CategoryType` via `category_type_id` (nullable — legacy rows may have none)
- Has many `Product` via `products.category_id`

**DTO:** `CreateCategoryDTO { Name, ImageURL string; CategoryTypeID uint }` — same file. When `CategoryTypeID` is sent, the category's `Name` is overwritten with the linked `CategoryType.Name` (see `internal/handlers/category.go`).

---

## 3. `CategoryType`
**Defined in:** `internal/models/base.go`
**Table:** `category_types` (default pluralization)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| Name | `string` | `unique;not null;type:varchar(255)` |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |
| DeletedAt | `gorm.DeletedAt` | `index` (soft delete) |

**Relations:** referenced by `Category.CategoryTypeID` (inverse of the belongs-to above).

**DTO:** `CreateCategoryTypeDTO { Name string }` — same file.

**Behavior note:** renaming a `CategoryType` (`PUT /api/categoryTypes/updateCategoryType/:id`) cascades — every `Category` row with that `category_type_id` gets its `name` bulk-updated to match, inside a transaction (`internal/handlers/category_type.go`).

---

## 4. `Product`
**Defined in:** `internal/models/product.go`
**Table:** `products` (default pluralization)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| Name | `string` | `not null;type:varchar(255)` |
| Brand | `string` | `type:varchar(100);default:'FITstore'` |
| Description | `string` | `type:text` |
| ProductCode | `string` | `unique;not null;type:varchar(100)` |
| SKU | `string` | `unique;not null;type:varchar(100)` |
| MRP | `float64` | `not null;type:decimal(10,2)` |
| SellingPrice | `float64` | `not null;type:decimal(10,2)` |
| Discount | `int` | `default:0` (column `discount`, JSON key `discount_percentage`) |
| Images | `string` | `type:text` — JSON/comma-serialized array of URLs |
| Specifications | `string` | `type:jsonb` |
| CategoryID | `uint` | `not null;index` |
| Category | `*Category` | `constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;` (belongs-to) |
| Sizes | `[]ProductSize` | `foreignKey:ProductID;constraint:OnDelete:CASCADE;` (has-many) |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |
| DeletedAt | `gorm.DeletedAt` | `index` (soft delete) |

**Relations:**
- Belongs to `Category` via `category_id`, DB-level `ON DELETE RESTRICT`. The app works around this at the handler level: `DeleteCategory` hard-deletes any linked products (and their sizes) first, inside the same transaction, before deleting the category.
- Has many `ProductSize` via `product_sizes.product_id`, DB-level `ON DELETE CASCADE`.

---

## 5. `ProductSize`
**Defined in:** `internal/models/product.go`
**Table:** `product_sizes` (default pluralization)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| ProductID | `uint` | `not null;index` |
| Size | `string` | `type:varchar(10);not null` |
| Stock | `int` | `not null;default:0` |
| Product | `*Product` | `foreignKey:ProductID;references:ID` (belongs-to) |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |
| DeletedAt | `gorm.DeletedAt` | `index` (soft delete) |

**Relations:** belongs to `Product` via `product_id`.

---

## 6. `AdminContent`
**Defined in:** `internal/content/model.go`
**Table:** `admin_content` — **explicit `TableName()` override** (default pluralization would have been `admin_contents`)

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint` | `primaryKey` |
| AboutUsImg | `string` | `type:text` |
| AboutUsTitle | `string` | `type:varchar(255)` |
| AboutUsDescription | `string` | `type:text` |
| AboutUsTagline1 | `string` | `type:varchar(255)` |
| AboutUsTagline2 | `string` | `type:varchar(255)` |
| AboutUsTagline3 | `string` | `type:varchar(255)` |
| AboutUsTagline4 | `string` | `type:varchar(255)` |
| AboutUsEstbYear | `string` | `type:varchar(10)` |
| AboutUsVisitUs | `string` | `type:text` |
| AboutUsEmail | `string` | `type:varchar(255)` |
| CreatedAt | `time.Time` | *(none)* |
| UpdatedAt | `time.Time` | *(none)* |

No `DeletedAt` (hard rows only). No relations. Only one row is ever expected to exist (About Us content).

**DTO:** `AboutUsDTO` mirrors every field above except `ID`/timestamps — same file.

---

## 7. `Order` — ⚠️ NOT owned by this service
**Defined in:** `internal/orders/model.go`
**Table:** `orders` — explicit `TableName()` override

> Owned and migrated by **fitstore-core (Prisma)**. Never pass this to `AutoMigrate` here — it's a read/update mapping onto a table this Go service doesn't create.

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint64` | `column:id;primaryKey` |
| CustUserID | `uint64` | `column:cust_user_id` |
| Status | `string` | `column:status` |
| ShippingName | `*string` | `column:shipping_name` |
| ShippingEmail | `*string` | `column:shipping_email` |
| ShippingAddress | `*string` | `column:shipping_address` |
| ShippingCity | `*string` | `column:shipping_city` |
| ShippingState | `*string` | `column:shipping_state` |
| ShippingPincode | `*string` | `column:shipping_pincode` |
| CreatedAt | `time.Time` | `column:created_at` |
| UpdatedAt | `time.Time` | `column:updated_at` |
| DecidedAt | `*time.Time` | `column:decided_at` |
| AdminComment | `*string` | `column:admin_comment` |
| Items | `[]OrderItem` | `foreignKey:OrderID` (has-many) |

⚠️ **Known drift:** `AdminComment` → `admin_comment` does not exist yet in fitstore-core's real Prisma schema (per the source comment). Reject-with-comment won't persist correctly until Prisma adds that column.

---

## 8. `OrderItem` — ⚠️ NOT owned by this service
**Defined in:** `internal/orders/model.go`
**Table:** `order_items` — explicit `TableName()` override

| Field | Go type | GORM tag |
|---|---|---|
| ID | `uint64` | `column:id;primaryKey` |
| OrderID | `uint64` | `column:order_id` |
| ProductSizeID | `uint64` | `column:product_size_id` |
| Quantity | `uint64` | `column:quantity` |
| UnitPrice | `float64` | `column:unit_price` |
| Status | `string` | `column:status` |
| CreatedAt | `time.Time` | `column:created_at` |
| UpdatedAt | `time.Time` | `column:updated_at` |
| ProductSize | `*models.ProductSize` | `foreignKey:ProductSizeID;references:ID` (belongs-to, cross-package into `internal/models`) |

`ProductSize`/`Product` are joined in purely for display.

---

## Migration ownership

`cmd/api/main.go` no longer calls `AutoMigrate`. This Go service only issues
plain reads/writes (`SELECT`/`INSERT`/`UPDATE`/`DELETE`) against tables it
assumes already exist with the exact columns documented above — it never
creates or alters schema itself, for any table.

**Practical consequence:** any future column/table addition on the Go side
(the way `category_type_id` was added to `categories` earlier) now requires a
matching Prisma migration in fitstore-core *first*, applied to the shared DB,
before the corresponding GORM struct/handler change here will work.

---

## Quick reference

| Model | Table | Owned/migrated by |
|---|---|---|
| `User` | `users` | Prisma (fitstore-core) |
| `Category` | `categories` | Prisma (fitstore-core) |
| `CategoryType` | `category_types` | Prisma (fitstore-core) |
| `Product` | `products` | Prisma (fitstore-core) |
| `ProductSize` | `product_sizes` | Prisma (fitstore-core) |
| `AdminContent` | `admin_content` | Prisma (fitstore-core) |
| `Order` | `orders` | Prisma (fitstore-core) |
| `OrderItem` | `order_items` | Prisma (fitstore-core) |
