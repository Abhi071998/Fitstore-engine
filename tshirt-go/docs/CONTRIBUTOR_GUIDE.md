# Contributor Guide — FitStore Admin Backend (tshirt-go)

Functional/technical reference for anyone writing code in this repo. For the
raw table/column reference, see [SCHEMA.md](SCHEMA.md). For a plain-language
explanation of what this system does, see [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md).

---

## 1. What this service is

`tshirt-go` is the **admin backend** for FitStore — a Go/Echo/GORM API that
powers the admin panel (a separate frontend, `fitstore-ui`). It manages
product catalog data (categories, category types, products) and admin-only
content (About Us page), and reviews/approves customer orders that are placed
through a **separate** customer-facing service, `fitstore-core` (Node.js +
Prisma), which shares the same PostgreSQL database.

**Ownership boundary (important):** as of the removal of `AutoMigrate`, this
Go service does not create or alter *any* database table, including the ones
it reads/writes most (`categories`, `products`, etc.). All schema/migrations
for the entire shared database are owned by fitstore-core's Prisma setup.
This service only issues plain `SELECT`/`INSERT`/`UPDATE`/`DELETE` queries
against tables that must already exist. See [SCHEMA.md](SCHEMA.md) for the
full column-level spec every model here assumes.

## 2. Tech stack

- **Language:** Go
- **Web framework:** [Echo v4](https://echo.labstack.com/)
- **ORM:** [GORM](https://gorm.io/) with the `gorm.io/driver/postgres` (pgx) dialector
- **Database:** PostgreSQL (shared with fitstore-core)
- **Auth:** JWT (HS256), via `github.com/golang-jwt/jwt/v5`
- **Deployment:** Render (`fitstore-engine.onrender.com`)

## 3. Project layout

```
cmd/api/main.go              — entrypoint: config load, DB connect, middleware, route wiring
internal/config/config.go    — env var loading (DATABASE_URL, PORT, JWT_SECRET, FRONTEND_URL, ENV)
internal/db/postgres.go      — GORM/pgx connection setup
internal/models/             — Category, CategoryType, Product, ProductSize, User models + DTOs
internal/handlers/           — HTTP handlers for auth, categories, category types, products
internal/content/            — About Us content model + handler (self-contained package)
internal/orders/             — Order/OrderItem models + handler (reads fitstore-core's tables)
internal/middleware/         — JWT auth middleware
internal/routes/routes.go    — single SetupRoutes() function wiring every route
docs/                        — SCHEMA.md, this file, PROJECT_OVERVIEW.md
```

There is no repository/service layer — handlers talk to `*gorm.DB` directly.
There is no dedicated migrations directory (see §1 — schema is Prisma-owned).

## 4. Local setup

```bash
git clone https://github.com/Abhi071998/Tshirt-admin-backend.git
cd tshirt-go
cp .env.example .env   # fill in DATABASE_URL etc.
go mod tidy
go run cmd/api/main.go
```

For faster iteration than `go run` (which recompiles every time):

```bash
go build -o server.exe ./cmd/api   # rebuild only after source changes
./server.exe                       # fast, no recompile
```

Server listens on `:8080` by default (`PORT` env var). Health check: `GET /health`.

### Required environment variables

| Var | Purpose |
|---|---|
| `DATABASE_URL` | Postgres connection string (shared DB — same one fitstore-core writes to) |
| `PORT` | defaults to `8080` |
| `ENV` | defaults to `development` |
| `JWT_SECRET` | HMAC signing key for JWTs; falls back to an insecure dev default if unset — **always set explicitly in production** |
| `FRONTEND_URL` | comma-separated list of allowed CORS origins, e.g. `http://localhost:5173,https://fitstore-ui.onrender.com` (no trailing slashes) |

## 5. Auth flow

1. `POST /api/auth/signup` — `{email, name, password}` → creates a `User` (bcrypt-hashed password).
2. `POST /api/auth/login` — `{email, password}` → returns `{token, user}`. Token is a 72-hour HS256 JWT carrying `user_id`/`email`.
3. Every protected route requires header `Authorization: Bearer <token>`, enforced by `middleware.JWTMiddleware` ([auth_middleware.go](../internal/middleware/auth_middleware.go)), which sets `user_id`/`user_email` on the Echo context for handlers to read.

## 6. Routes

🌐 public · 🔒 requires Bearer token. Full wiring in [routes.go](../internal/routes/routes.go).

| Method | Route | Access | Notes |
|---|---|---|---|
| GET | `/health` | 🌐 | liveness/DB-sync check |
| POST | `/api/auth/signup` | 🌐 | |
| POST | `/api/auth/login` | 🌐 | returns JWT |
| GET | `/api/categoryTypes/getAllCategoryTypes` | 🔒 | fixed dropdown list |
| POST | `/api/categoryTypes/createCategoryType` | 🔒 | |
| PUT | `/api/categoryTypes/updateCategoryType/:id` | 🔒 | **cascades** rename to every linked `Category.name` |
| DELETE | `/api/categoryTypes/deleteCategoryType/:id` | 🔒 | |
| GET | `/api/categories/getAllCategories` | 🔒 | preloads `Products` |
| POST | `/api/categories/createCategory` | 🔒 | see §7 |
| PUT | `/api/categories/updateCategory/:id` | 🔒 | see §7 |
| DELETE | `/api/categories/deleteCategory/:id` | 🔒 | **cascades**: hard-deletes linked products + their sizes first, see §7 |
| GET | `/api/products/getAllProducts/:categoryId` | 🔒 | preloads `Category`, `Sizes` |
| POST | `/api/products/createProduct` | 🔒 | |
| PUT | `/api/products/updateProduct/:id` | 🔒 | partial update via `map[string]interface{}`; replaces all `Sizes` if `sizes` is sent |
| DELETE | `/api/products/deleteProduct/:id` | 🔒 | hard delete (product + sizes), not soft |
| GET | `/api/orders/pending` | 🔒 | grouped by customer; returns `[]` (200) instead of 500 if the `orders` table doesn't exist yet — see §8 |
| PUT | `/api/orders/:id/approve` | 🔒 | |
| PUT | `/api/orders/:id/reject` | 🔒 | requires `comment` in body |
| POST | `/api/orders/bulk-approve` | 🔒 | `{order_ids: []}`, silently skips non-pending/missing IDs |
| GET | `/api/content/about-us` | 🌐 | 404 until the single row is created |
| POST | `/api/content/about-us` | 🔒 | 409 if a row already exists — only one row is ever expected |
| PUT | `/api/content/about-us` | 🔒 | |

## 7. Category / CategoryType relationship

`CategoryType` ([base.go](../internal/models/base.go)) is a fixed, admin-managed
dropdown list (`{id, name}`). `Category` optionally links to one via a nullable
`CategoryTypeID` FK.

- **Create/update a category**: if `category_type_id` is sent in the payload,
  the handler looks it up and the category's `name` is **overwritten** with
  the type's name (the dropdown value wins over any free-typed `name`). If
  `category_type_id` is omitted/zero, the legacy free-text `name` still works.
- **Rename a category type**: `UpdateCategoryType` runs a transaction that
  renames the type, then bulk-updates `name` on every `Category` row with
  that `category_type_id` — see [category_type.go](../internal/handlers/category_type.go).
- **Delete a category**: `DeleteCategory` ([category.go](../internal/handlers/category.go))
  first finds every `Product` under it, hard-deletes their `ProductSize` rows,
  then hard-deletes the products, then hard-deletes the category itself —
  all in one transaction. The response reports `deleted_products: N`.

## 8. Orders module — cross-service dependency

`internal/orders` reads/writes the `orders`/`order_items` tables, which are
owned and created by **fitstore-core** (Prisma), not this service. Until
fitstore-core is deployed and has run its own migrations, those tables don't
exist yet in a fresh database, which previously caused a hard 500
(`relation "orders" does not exist`, SQLSTATE `42P01`) on `GET /api/orders/pending`.

That endpoint now catches Postgres error `42P01` specifically and returns
`200 []` instead, so the admin dashboard degrades gracefully instead of
erroring while fitstore-core isn't live yet. `approve`/`reject`/`bulk-approve`
are not guarded this way (nothing to act on if the list is empty).

**Known schema drift:** `Order.AdminComment` maps to `admin_comment`, which
per a source comment does not yet exist in fitstore-core's real Prisma
schema — reject-with-comment won't persist correctly until that column is
added on the Prisma side.

## 9. Known infra gotchas

- **CORS**: `FRONTEND_URL` (comma-separated) drives `AllowOrigins` in
  [main.go](../cmd/api/main.go). A frontend origin missing from this list
  fails as a CORS-blocked preflight, which can visually mask an underlying
  4xx/5xx from the actual endpoint in the browser console — always check
  the raw response status in the Network tab, not just the console error text.
- **"cached plan must not change result type" (SQLSTATE 0A000)**: happens
  when a live deploy alters a table's columns while existing pgx connections
  still hold prepared plans against the old shape. Fixed by opening the GORM
  postgres dialector with `PreferSimpleProtocol: true` ([postgres.go](../internal/db/postgres.go)),
  which disables prepared-statement caching entirely. If this ever recurs
  after a schema change made outside this fix, a Render service restart
  clears it immediately.
- **Soft vs hard delete**: `Category`/`Product`/`User`/`CategoryType` all
  carry `gorm.DeletedAt`, so a plain `.Delete()` is a *soft* delete (sets
  `deleted_at`, row stays in the table, hidden from normal queries). Product
  and category deletion in this codebase deliberately use `.Unscoped().Delete()`
  for a real, permanent removal — don't assume GORM's default `.Delete()`
  behavior when reading/writing delete logic here.
