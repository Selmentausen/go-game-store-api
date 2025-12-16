# 🎮 Gopher Game Store API

A production-ready E-Commerce System built with **Go (Golang)**, featuring Microservices (gRPC), Clean Architecture, and a responsive Frontend.

![Go](https://img.shields.io/badge/Go-1.25-blue) ![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2B%20Microservices-green) ![Docker](https://img.shields.io/badge/Docker-Enabled-blue)

## 🏗 Architecture
This project demonstrates a transition from a Monolith to a **Microservices** architecture.

*   **Game Store API (REST):** The main entry point. Handles Users, Products, and Cart logic.
*   **Payment Service (gRPC):** A separate, isolated service that processes payments. The API talks to this service via Protocol Buffers.
*   **Database:** PostgreSQL (Relation data) and Redis (Async Job Queue).
*   **Pattern:** Logic is strictly separated into `Handlers`, `Services`, and `Repositories`.

## 🚀 Tech Stack
*   **Language:** Go (Golang)
*   **Communication:** REST (Gin) & gRPC (Protobuf)
*   **Database:** PostgreSQL 16 & Redis
*   **Frontend:** Vanilla JS + Tailwind CSS (SPA)
*   **Testing:** SQLite (In-Memory) & Testify
*   **Orchestration:** Docker Compose

## 📂 Project Structure
```text
/cmd
  /api           # Main REST API Server
  /seeder        # Data population tool
/internal
  /handlers      # HTTP Controllers
  /service       # Business Logic
  /repository    # Data Access (GORM)
  /grpc          # Generated Protobuf code
/payment-service # Microservice
  /cmd/server    # gRPC Server Entry
  /proto         # Protocol Buffer Definitions
/static          # Frontend Assets
```

## 🛠 Installation & Usage

### 1. Run with Docker Compose (Recommended)
This command spins up the API, Payment Service, Postgres, and Redis in a private network.
```bash
docker-compose up --build
```
*   **Frontend:** Visit `http://localhost:8080`
*   **API:** `http://localhost:8080/api/v1/...`

### 2. Seed the Database
Populate the store with dummy data (must run against the exposed Docker ports).
```bash
go run cmd/seeder/main.go
```
*   **Admin Login:** `admin@gamestore.com` / `password123`
*   **User Login:** `player1@test.com` / `password123`

## 🛒 Features Implemented

### Microservice Payments
*   When `Checkout` is triggered, the Order Service creates a gRPC connection to the Payment Service.
*   The Payment Service validates limits and returns a transaction ID.
*   The Order is only saved if the gRPC call returns `Success: true`.

### 3. Concurrency & Async
*   **Job Queue:** Registration triggers a "Welcome Email" task pushed to Redis.
*   **Worker Pool:** A background goroutine consumes tasks from Redis to prevent blocking the API.

## 🔑 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **Auth** | | |
| POST | `/api/v1/auth/login` | Get JWT Token |
| **Cart** | | |
| GET | `/api/v1/cart` | View Cart |
| POST | `/api/v1/cart` | Add/Update Item (qty: 1 or -1) |
| DELETE | `/api/v1/cart/:id` | Remove Item completely |
| POST | `/api/v1/cart/checkout` | Process Payment & Order |
| **Products** | | |
| GET | `/api/v1/products` | List Inventory |