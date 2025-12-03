# 🎮 Gopher Game Store API

A high-performance, concurrent E-Commerce Backend built with **Go (Golang)**.
This project demonstrates production-ready patterns including **Database Transactions**, **Asynchronous Task Queues**, **RBAC Security**, and **Graceful Shutdowns**.

![Go](https://img.shields.io/badge/Go-1.23-blue) ![Docker](https://img.shields.io/badge/Docker-Enabled-blue) ![Redis](https://img.shields.io/badge/Redis-Async-red)

## 🏗 Architecture
The application follows a Clean Architecture approach separating Handlers, Models, and Middleware.

*   **Core API:** Built with **Gin** for high-performance HTTP routing.
*   **Database:** **PostgreSQL** managed via **GORM** (Code-first migrations).
*   **Concurrency:** Heavy tasks (Email simulation) are offloaded to **Redis** and processed by background **Goroutines**.
*   **Data Integrity:** Purchasing logic uses **ACID Transactions** with row-level locking to prevent race conditions during high traffic.
*   **Security:** JWT Authentication with Role-Based Access Control (RBAC).

## 🚀 Tech Stack
*   **Language:** Go (Golang)
*   **Framework:** Gin Web Framework
*   **Databases:** PostgreSQL 16 (Primary), Redis (Queue/Cache)
*   **Containerization:** Docker & Docker Compose
*   **Testing:** SQLite (In-Memory) for Unit/Integration tests

## 🛠 Installation & Setup

### Prerequisites
*   Go 1.21+
*   Docker Desktop

### 1. Start Infrastructure
Spin up PostgreSQL and Redis containers.
```bash
docker-compose up -d
```

### 2. Run the Server
```bash
go run main.go
```
The server will start on `http://localhost:8080`.

### 3. Create Admin User (Seeding)
Registration creates standard users by default. To create an Admin:
```bash
go run cmd/admin-seeder/main.go
```
*Credentials: admin@gamestore.com / admin123*

## ✅ Running Tests
Tests use an in-memory SQLite database to mock the Postgres connection, allowing for fast, isolated execution.
```bash
go test ./handlers -v
```

## 🔑 API Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/register` | Create a new user account | No |
| POST | `/api/v1/auth/login` | Login and receive JWT | No |
| GET | `/api/v1/products` | View inventory | No |
| POST | `/api/v1/products` | Add new product (Admin only) | **Yes (Admin)** |
| POST | `/api/v1/orders` | Purchase a game (Transactional) | **Yes** |

## 🧪 Key Features Implementation

### 1. Database Transactions (Safety)
Located in `handlers/order.go`.
Ensures that Inventory is only deducted if the Order is successfully created. Uses `FOR UPDATE` locking to prevent two users buying the last item simultaneously.

### 2. Async Workers (Performance)
Located in `worker/email_worker.go`.
Registration triggers a "Welcome Email" task pushed to a Redis List. A dedicated Goroutine consumes this list, preventing the API from blocking during slow I/O operations.

### 3. Graceful Shutdown (Reliability)
Located in `main.go`.
The server listens for OS Interrupt signals (SIGINT/SIGTERM). When received, it allows active requests 5 seconds to complete before closing connections.