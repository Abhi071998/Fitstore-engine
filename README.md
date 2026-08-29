# 👕 FitStore-engine (Admin-backend)

A lightweight administrative backend engine built with **Go (Golang)**, the **Echo web framework**, and **GORM ORM** using a **PostgreSQL** database.

---

## 🛠️ Prerequisites

- Go 1.21+
- PostgreSQL 14+
- [TablePlus](https://tableplus.com/) (or any Postgres GUI/CLI) for initial DB creation

---

## 📥 Clone the Project

```bash
git clone https://github.com/Abhi071998/Tshirt-admin-backend.git
cd tshirt-go
```

## 🐘 Database Setup (TablePlus)

Before running the backend, create the database container inside PostgreSQL:

1. Open TablePlus and connect to your local PostgreSQL server instance.
2. Press `Ctrl + G` (or `Cmd + G` on Mac) to open the database panel.
3. Click **"New..."** or **"+"**, type the database name exactly as `tshirt_store`, and click Save/Create.

> **Note:** You only need to create the empty database. All tables (including the new `category_types` table) are created automatically on server startup via GORM `AutoMigrate` — no manual migration step is required.

## 🔑 Environment Setup

Create a `.env` file in the root directory and populate it from `.env.example`:

```env
PORT=8080
ENV=local
JWT_SECRET=your_placeholder_secret_key
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/tshirt_store?sslmode=disable
FRONTEND_URL=http://localhost:5173
```

## 🚀 Run the Application

From the `tshirt-go` directory:

```bash
# Clean up and download required packages
go mod tidy

# Start the server (recompiles from source every time — slower)
go run cmd/api/main.go
```

### ⚡ Faster local runs (build once, reuse the binary)

`go run` recompiles the whole binary on every start. For faster iteration, build it once and just re-run the resulting executable:

```bash
# Build once
go build -o server.exe ./cmd/api

# Then just run the binary directly (fast — no recompile)
./server.exe
```

Re-run `go build -o server.exe ./cmd/api` only after you change Go source code; otherwise just re-launch `./server.exe`.

The server starts on `http://localhost:8080` by default (configurable via `PORT`). You can verify it's up by hitting `GET /health`.
