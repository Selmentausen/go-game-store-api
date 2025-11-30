# Go Game Store API

A high-performance REST API built with Go (Golang), Gin, PostgreSQL, and Redis.
Implements asynchronous task processing for email notifications.

## Tech Stack
- **Language:** Golang
- **Framework:** Gin
- **Database:** PostgreSQL (GORM)
- **Queue/Cache:** Redis
- **Auth:** JWT & Bcrypt
- **Infrastructure:** Docker Compose

## Features
- JWT Authentication (Login/Register)
- Product CRUD Operations
- **Concurrency:** Background worker for email simulation using Redis Lists
- **Security:** Password Hashing, Middleware protection
- **Reliability:** Graceful Shutdown & Dockerized environment

## How to Run
1. Clone the repo
2. `docker-compose up -d`
3. `go run main.go`