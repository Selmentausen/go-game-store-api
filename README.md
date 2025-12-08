# 🎮 Gopher Game Store API

A production-ready E-Commerce Backend built with **Go (Golang)**, featuring Clean Architecture, ACID transactions, and asynchronous processing. Includes a responsive Frontend for demonstration.

![Go](https://img.shields.io/badge/Go-1.23-blue) ![Architecture](https://img.shields.io/badge/Architecture-Clean-green) ![Docker](https://img.shields.io/badge/Docker-Enabled-blue) ![Redis](https://img.shields.io/badge/Redis-Async-red)

## 🏗 Architecture & Design Patterns

*   **Internal Layering:** Logic is strictly separated into `Handlers` (HTTP), `Services` (Business Logic), and `Repositories` (Data Access).
*   **Dependency Injection:** All layers are injected via constructors in `main.go`, making the system loosely coupled and highly testable.
*   **Atomic Transactions:** The Checkout process uses database transactions to ensure Inventory and Orders stay consistent even during failures.
*   **Concurrency:** Heavy tasks (Email simulation) are offloaded to **Redis** and processed by background **Goroutines**.
*   **Structured Logging:** Uses Go's `log/slog` for JSON-formatted production logs.

## 🚀 Tech Stack
*   **Language:** Go (Golang)
*   **Web Framework:** Gin
*   **Database:** PostgreSQL 16 (Managed via GORM)
*   **Caching/Queue:** Redis
*   **Frontend:** HTML5, Vanilla JS, Tailwind CSS
*   **Testing:** SQLite (In-Memory)
*   **Containerization:** Docker & Docker Compose

## 📂 Project Structure
```text
/cmd
  /api           # Application Entry Point (Composition Root)
  /seeder        # Data Population Scripts
/internal
  /models        # Domain Entities
  /handlers      # HTTP Transport Layer
  /service       # Business Logic Layer
  /repository    # Database Access Layer
  /worker        # Background Job Processors
  /middleware    # JWT Auth & RBAC
/static          # Frontend Assets
```

## 🛠 Installation & Setup

### Prerequisites
*   Go 1.21+
*   Docker Desktop

### 1. Start Infrastructure
Spin up PostgreSQL and Redis containers.
```bash
docker-compose up -d
```

### 2. Seed the Database
Populates the database with demo products and users.
```bash
go run cmd/seeder/main.go
```
*   **Admin:** `admin@gamestore.com` / `password123`
*   **User:** `player1@test.com` / `password123`

### 3. Run the Server
```bash
go run cmd/api/main.go
```
The server will start on **http://localhost:8080**.

## 🛒 Features

### 1. Shopping Cart & Checkout
*   Full cart management (Add, Remove, View).
*   **Transactional Checkout:** `internal/service/order_service.go` performs a complex transaction:
    1.  Locks product rows (`FOR UPDATE`) to prevent race conditions.
    2.  Checks stock levels.
    3.  Deducts stock.
    4.  Creates Order header and Order Items.
    5.  Clears Cart.
    6.  Commits transaction only if ALL steps succeed.

### 2. Authentication & RBAC
*   **JWT** based stateless authentication.
*   **Middleware** protection for routes (`AuthMiddleware`).
*   **Role-Based Access Control:** Only Admins can add products (`AdminOnly` middleware).

### 3. Asynchronous Workers
*   Registration triggers a "Welcome Email" task.
*   The API pushes this task to a **Redis List** (Queue).
*   A background worker (`internal/worker`) consumes the task and simulates sending an email without blocking the HTTP response.

### 4. Frontend
*   A responsive Single Page Application (SPA) located at `/static`.
*   Uses **Fetch API** to communicate with the Go backend.
*   Includes dynamic Cart UI, Toast notifications, and Admin Dashboard.

## 🧪 Testing
Unit and Integration tests run against an **In-Memory SQLite** database to mock PostgreSQL, ensuring fast and isolated execution.

```bash
go test ./internal/handlers -v
```

## 🔑 API Endpoints

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/register` | Create account | No |
| POST | `/api/v1/auth/login` | Get JWT Token | No |
| GET | `/api/v1/products` | List Inventory | No |
| POST | `/api/v1/products` | Add Product | **Admin** |
| GET | `/api/v1/cart` | View Cart | **User** |
| POST | `/api/v1/cart` | Add Item to Cart | **User** |
| POST | `/api/v1/cart/checkout` | Process Order | **User** |