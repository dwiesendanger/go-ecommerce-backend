# Scalable E-Commerce Backend in Go

A high-performance, scalable RESTful API for an e-commerce platform built with **Go (Golang)**. This project demonstrates modern backend practices including Clean Architecture, asynchronous processing with Worker Pools, Redis caching, and containerization.

## Features

* **Clean Architecture:** Clear separation of concerns (Handlers, Services, Domain Models).
* **High Performance:** Built on top of [Gin](https://github.com/gin-gonic/gin), one of the fastest Go web frameworks.
* **Database:** PostgreSQL with **GORM** for robust data management.
* **Caching:** **Redis** implementation (Cache-Aside pattern) for product listings to minimize DB load.
* **Concurrency:** Asynchronous **Worker Pool** pattern using Goroutines and Channels to handle background tasks (e.g., simulating email sending and inventory updates) without blocking the API.
* **Authentication:** Stateless JWT (JSON Web Token) authentication with Bcrypt password hashing.
* **ACID Transactions:** Safe order processing ensuring inventory and order data integrity.
* **Graceful Shutdown:** Handles `SIGINT`/`SIGTERM` to ensure background workers finish tasks before the server stops.
* **Containerization:** Full **Docker** and **Docker Compose** support for easy deployment.

## Tech Stack

* **Language:** Go 1.23+
* **Framework:** Gin Web Framework
* **Database:** PostgreSQL 15
* **Cache:** Redis
* **ORM:** GORM
* **Auth:** JWT (golang-jwt) & Bcrypt
* **Deployment:** Docker & Docker Compose

## Project Structure

```text
ecommerce-platform/
├── cmd/
│   └── api/
│       └── main.go        # Application entry point
├── internal/
│   ├── core/
│   │   └── domain/        # Database models (User, Product, Order, Cart)
│   ├── handlers/          # HTTP Controllers (Input validation, JSON response)
│   ├── services/          # Business Logic (Caching, Transactions)
│   └── workers/           # Background Worker Pool (Async tasks)
├── pkg/
│   ├── auth/              # JWT Token generation & validation
│   └── database/          # Postgres & Redis connection setup
├── Dockerfile             # Multi-stage build definition
├── docker-compose.yml     # Orchestration for App, DB, and Redis
└── go.mod                 # Dependencies
```

## Getting Started

The easiest way to run the application is using Docker Compose.

### Prerequisites

* Docker
* Docker Compose

### Running with Docker (Recommended)

1. **Clone the repository:**
```bash
git clone [https://github.com/dwiesendanger/go-ecommerce-backend.git](https://github.com/dwiesendanger/go-ecommerce-backend.git)
cd go-ecommerce-backend
```


2. **Start the application:**
```bash
docker-compose up --build
```

*This will start PostgreSQL, Redis, and the Go API.*

3. **Access the API:**
The server will start on `http://localhost:8080`.

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| POST | `/api/v1/register` | Register a new user | No |
| POST | `/api/v1/login` | Login and receive JWT Token | No |

### Products

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| GET | `/api/v1/products` | List all products (Cached) | No |
| POST | `/api/v1/products` | Create a new product (Admin) | **Yes** |

### Cart

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| GET | `/api/v1/cart` | View current user's cart | **Yes** |
| POST | `/api/v1/cart/items` | Add item to cart | **Yes** |

### Orders

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| POST | `/api/v1/orders` | Checkout (Convert cart to order) | **Yes** |
| GET | `/api/v1/orders` | View order history | **Yes** |

## Testing the Flow

1. **Register:** `POST /register` with email/password.
2. **Login:** `POST /login` to get the `token`.
3. **Create Product:** `POST /products` (use the token in `Authorization: Bearer <token>` header).
4. **Add to Cart:** `POST /cart/items` with `product_id`.
5. **Checkout:** `POST /orders`.
* *Observation:* The API responds immediately. Check the server logs to see the **Background Workers** processing the email and inventory tasks asynchronously.

6. **History:** `GET /orders` to see your past orders.

## License

Distributed under the MIT License. See `LICENSE` for more information.