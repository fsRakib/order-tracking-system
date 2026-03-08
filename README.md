# Order Tracking System

> **Learning Project:** A hands-on microservices application built to learn and demonstrate modern backend technologies including **gRPC**, **Elasticsearch**, **RabbitMQ**, **PostgreSQL**, and **Microservice Architecture**.

## 🎯 Project Purpose

This project serves as a practical learning platform to understand:

- **gRPC** - High-performance RPC framework for inter-service communication
- **Elasticsearch** - Distributed search and analytics engine
- **RabbitMQ** - Message broker for event-driven architecture
- **PostgreSQL** - Relational database management
- **Microservices** - Distributed system design and implementation
- **Domain-Driven Design (DDD)** - Structuring code around business domains with aggregates, value objects, domain events, and the ports-and-adapters pattern

---

## 📊 Technology Architecture & Usage

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                    🚀 ORDER TRACKING SYSTEM ARCHITECTURE                     ║
╚══════════════════════════════════════════════════════════════════════════════╝

                            👥 CLIENT LAYER
    ┌──────────────────────────────────────────────────────────┐
    │  🔌 gRPC Clients          │    🌐 HTTP/REST Clients      │
    │  (Order & Stock Services) │    (Analytics API)           │
    └──────────┬────────────────┴─────────────┬────────────────┘
               │                               │
               │                               │
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                        APPLICATION LAYER (Microservices)
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
               │                               │
    ┏━━━━━━━━━▼━━━━━━━━━━┓         ┏━━━━━━━━▼━━━━━━━━━━━━━━━━┓
    ┃  📦 ORDER SERVICE  ┃         ┃  📊 ANALYTICS SERVICE   ┃
    ┃  ═════════════════ ┃         ┃  ══════════════════════ ┃
    ┃  🛠️  gRPC Server    ┃         ┃  🌐 HTTP REST API      ┃
    ┃  📍 Port: 50061    ┃         ┃  📍 Port: 8081          ┃
    ┃                    ┃         ┃                         ┃
    ┃  Operations:       ┃         ┃  Endpoints:             ┃
    ┃  • CreateOrder     ┃         ┃  • GET /search          ┃
    ┃  • GetOrder        ┃         ┃  • GET /aggregate/...   ┃
    ┃  • UpdateStatus    ┃         ┃  • GET /health          ┃
    ┃  • GetByCustomer   ┃         ┃                         ┃
    ┗━━━━━┯━━━━━━┯━━━━━━━┛         ┗━━━━━━━━━━┯━━━━━━━━━━━━━━┛
          │      │                           │
          │      │ ⚡ gRPC Call              │
          │      │ (Stock Check)             │
          │      │                           │
          │  ┏━━━▼━━━━━━━━━━━━━┓             │
          │  ┃ 📦 STOCK SERVICE ┃             │
          │  ┃ ════════════════ ┃             │
          │  ┃ 🛠️  gRPC Server  ┃             │
          │  ┃ 📍 Port: 50062   ┃             │
          │  ┃                  ┃             │
          │  ┃ Operations:      ┃             │
          │  ┃ • ReserveStock   ┃             │
          │  ┃ • ReleaseStock   ┃             │
          │  ┃ • GetStock       ┃             │
          │  ┗━━━━━━┯━━━━━━━━━━━┛             │
          │         │                        │
          │         │                        │
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
            MESSAGE BROKER LAYER (Event-Driven)
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
          │         │                        │
          │    ┏━━━━▼━━━━━━━━━━━━━━━━━━━━━━━▼━━━━━━━━━━━━┓
          │    ┃      🐰 RABBITMQ (Message Broker)       ┃
          │    ┃      ═══════════════════════════════    ┃
          │    ┃      📍 Port: 5673 (AMQP Protocol)      ┃
          │    ┃      🖥️  Management UI: 15673           ┃
          │    ┃                                         ┃
          │    ┃      📮 Queue: order_events             ┃
          │    ┃      ┌──────────────────────────┐       ┃
          │    ┃      │ 📨 OrderCreated          │       ┃
          └────┃──────│ 📨 OrderStatusUpdated    │       ┃
  Publisher    ┃      └──────────────────────────┘       ┃
               ┃           │              │              ┃
               ┗━━━━━━━━━━━┿━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━┛
                           │              │
                    Consumers (Async Processing)
                           │              │
                           ▼              ▼
                  ┌─── Stock Svc  Analytics Svc ───┐
                  │   (Reserve/      (Index to     │
                  │    Release)     Elasticsearch)  │
                  └─────────────────────────────────┘
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                    DATA PERSISTENCE LAYER
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
          │         │                        │
          │         │                        │
    ┏━━━━━▼━━━━━━━━━▼━━━━━━━┓      ┏━━━━━━━▼━━━━━━━━━━━━━━━━┓
    ┃  🗄️  POSTGRESQL DB    ┃      ┃  🔍 ELASTICSEARCH      ┃
    ┃  ══════════════════   ┃      ┃  ═════════════════     ┃
    ┃  📍 Port: 5433        ┃      ┃  📍 Port: 9201         ┃
    ┃  🗃️  Database: orders ┃      ┃  📑 Index: orders      ┃
    ┃                       ┃      ┃                        ┃
    ┃  📋 Tables:           ┃      ┃  ⚡ Capabilities:       ┃
    ┃  • customers          ┃      ┃  • Full-text search    ┃
    ┃  • orders             ┃      ┃  • Aggregations        ┃
    ┃  • order_items        ┃      ┃  • Real-time analytics ┃
    ┃  • stocks             ┃      ┃  • Fast queries        ┃
    ┃                       ┃      ┃                        ┃
    ┃  ✨ ACID Transactions ┃      ┃  ✨ Search & Analytics ┃
    ┗━━━━━━━━━━━━━━━━━━━━━━━┛      ┗━━━━━━━━━━━━━━━━━━━━━━━━┛

╔══════════════════════════════════════════════════════════════════════════════╗
║  🔄 DATA FLOW: Client → Service → RabbitMQ → Consumers → Database/ES         ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### � Request & Event Flow Visualization

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                         EXAMPLE: CREATE ORDER FLOW                           ║
╚══════════════════════════════════════════════════════════════════════════════╝

STEP 1: Client Request (Synchronous gRPC)
────────────────────────────────────────────────────────────────────────────────

    👤 Client
       │
       │ 1️⃣  gRPC: CreateOrder(customer_id, items[])
       ▼
    ┌─────────────────┐
    │ ORDER SERVICE   │  2️⃣  Validate request
    │   (Port 50061)  │  3️⃣  Generate order_id
    └────┬───────┬────┘  4️⃣  Calculate total_amount
         │       │
         │       │ 5️⃣  gRPC: CheckStock(sku, quantity)
         │       ▼
         │    ┌─────────────────┐
         │    │ STOCK SERVICE   │  6️⃣  Check available stock
         │    │   (Port 50062)  │  7️⃣  Return: success/failure
         │    └─────────────────┘
         │       │
         │       ▼ (if stock available)
         │
         │ 8️⃣  Save to PostgreSQL (orders, order_items tables)
         ▼
    ┌─────────────────┐
    │  POSTGRESQL DB  │
    └─────────────────┘
         │
         │ 9️⃣  Return Response to Client
         ▼
    👤 Client receives: {order_id, status, message}


STEP 2: Asynchronous Event Processing (Event-Driven)
────────────────────────────────────────────────────────────────────────────────

    ┌─────────────────┐
    │ ORDER SERVICE   │ 🔟 Publish event: OrderCreated
    └────────┬────────┘
             │
             │ Event: {order_id, customer_id, items, status, total}
             ▼
    ┌─────────────────────────────┐
    │  🐰 RABBITMQ                │
    │  Queue: order_events        │
    └────┬───────────────────┬────┘
         │                   │
         │ 1️⃣1️⃣ Fanout       │ 1️⃣2️⃣ Fanout
         │                   │
         ▼                   ▼
    ┌──────────────┐    ┌────────────────────┐
    │ STOCK SVC    │    │ ANALYTICS SVC      │
    │ Consumer     │    │ Consumer           │
    └──────┬───────┘    └────────┬───────────┘
           │                     │
           │ 1️⃣3️⃣ Reserve       │ 1️⃣4️⃣ Index order
           │    Stock            │    data
           ▼                     ▼
    ┌──────────────┐    ┌────────────────────┐
    │ POSTGRESQL   │    │ ELASTICSEARCH      │
    │ UPDATE stocks│    │ orders index       │
    │ quantity -= n│    │ {searchable data}  │
    └──────────────┘    └────────────────────┘


RESULT: Order Analytics Now Searchable!
────────────────────────────────────────────────────────────────────────────────

    👤 Another Client
       │
       │ GET /search?customer=John
       ▼
    ┌─────────────────────┐
    │ ANALYTICS SERVICE   │ Search Elasticsearch
    │   (Port 8080)       │ ⚡ Fast full-text search
    └──────────┬──────────┘
               │
               │ Query ES index
               ▼
    ┌─────────────────────┐
    │  ELASTICSEARCH      │ Returns matching orders
    └──────────┬──────────┘
               │
               ▼
    👤 Client receives: {total: 5, orders: [...]}


═══════════════════════════════════════════════════════════════════════════════
  KEY PATTERNS DEMONSTRATED:

  ✅ Synchronous Communication:  gRPC (Order ↔ Stock) - Fast, low latency
  ✅ Asynchronous Communication: RabbitMQ - Decoupled, resilient
  ✅ Event-Driven Architecture:  Publish/Subscribe pattern
  ✅ Data Consistency:           PostgreSQL transactions
  ✅ Search Performance:         Elasticsearch indexing
═══════════════════════════════════════════════════════════════════════════════
```

### �🔧 Technology Usage Breakdown

| Technology        | Used In                       | Purpose                                                                 |
| ----------------- | ----------------------------- | ----------------------------------------------------------------------- |
| **gRPC**          | Order Service ↔ Stock Service | Fast, type-safe inter-service communication with Protocol Buffers       |
| **PostgreSQL**    | Order Service, Stock Service  | Persistent storage for orders, customers, stocks with ACID transactions |
| **RabbitMQ**      | All Services                  | Event-driven communication, decoupling services, async processing       |
| **Elasticsearch** | Analytics Service             | Fast full-text search, aggregations, and real-time analytics on orders  |
| **HTTP REST**     | Analytics Service             | External API for search and analytics queries                           |
| **Microservices** | System Architecture           | Independent, scalable services with separate concerns                   |

---

## 📁 Project Structure

```
order-tracking-system/
│
├── order-service/              # Microservice #1: Order Management (DDD)
│   ├── main.go                 # Composition root - wires all dependencies
│   ├── Dockerfile              # Docker build configuration
│   ├── go.mod                  # Go module definition
│   │
│   ├── domain/                 # Core business logic - no external dependencies
│   │   ├── aggregate/
│   │   │   ├── order.go        # Order aggregate root (business rules)
│   │   │   └── order_item.go   # Order item entity
│   │   ├── valueobject/
│   │   │   ├── order_id.go     # UUID-based order identity
│   │   │   ├── customer_id.go  # Validated customer identifier
│   │   │   ├── order_status.go # Status enum with transition rules
│   │   │   ├── money.go        # Monetary value (stored as cents)
│   │   │   └── sku.go          # Stock keeping unit identifier
│   │   ├── event/
│   │   │   └── order_events.go # Domain events (OrderCreated, StatusUpdated)
│   │   └── repository/
│   │       └── order_repository.go  # Repository interface (port)
│   │
│   ├── application/            # Use cases - orchestrates domain objects
│   │   ├── command/
│   │   │   ├── create_order.go         # CreateOrder use case handler
│   │   │   └── update_order_status.go  # UpdateOrderStatus use case handler
│   │   ├── query/
│   │   │   └── order_queries.go  # GetOrder, GetOrdersByCustomer handlers
│   │   ├── dto/
│   │   │   └── order_dto.go      # Data transfer objects (layer boundary)
│   │   └── port/
│   │       ├── event_publisher.go  # EventPublisher interface
│   │       ├── stock_service.go    # StockService interface
│   │       └── customer_service.go # CustomerService interface
│   │
│   ├── infrastructure/         # Concrete implementations of interfaces
│   │   ├── persistence/
│   │   │   ├── postgres_order_repository.go    # SQL implementation
│   │   │   └── postgres_customer_service.go    # SQL implementation
│   │   ├── messaging/
│   │   │   └── rabbitmq_publisher.go  # RabbitMQ event publishing
│   │   └── grpc_client/
│   │       └── stock_grpc_client.go   # gRPC call to stock-service
│   │
│   ├── interface/              # Delivery layer - adapts external protocols
│   │   └── grpc/
│   │       └── order_handler.go  # Thin gRPC handler, delegates to application
│   │
│   └── db/
│       ├── db.go               # Database connection
│       └── schema.sql          # Database schema (customers, orders, order_items)
│
├── stock-service/              # Microservice #2: Stock Management
│   ├── main.go                 # Entry point
│   ├── Dockerfile              # Docker build configuration
│   ├── server/
│   │   └── server.go           # gRPC server (3 methods)
│   ├── consumer/
│   │   └── consumer.go         # Consumes order.cancelled events
│   └── db/
│       ├── db.go               # Database connection
│       ├── queries.go          # SQL queries
│       └── schema.sql          # stocks table schema
│
├── analytics-service/          # Microservice #3: Analytics & Search
│   ├── main.go                 # Entry point
│   ├── Dockerfile              # Docker build configuration
│   ├── api/
│   │   └── handler.go          # HTTP REST handlers (4 endpoints)
│   ├── consumer/
│   │   └── consumer.go         # Consumes order events, indexes to ES
│   └── elastic/
│       ├── client.go           # Elasticsearch connection
│       └── index.go            # Indexing operations
│
├── pb/                         # Shared generated Protocol Buffer code
│   ├── go.mod                  # Standalone Go module (order-tracking-system/pb)
│   ├── order/                  # Order service generated code
│   │   ├── order.pb.go
│   │   └── order_grpc.pb.go
│   └── stock/                  # Stock service generated code
│       ├── stock.pb.go
│       └── stock_grpc.pb.go
│
├── proto/                      # Protocol Buffer source definitions
│   ├── order.proto
│   └── stock.proto
│
├── go.work                     # Go workspace (links all 4 modules)
├── docker-compose.yml          # Docker orchestration (all services)
└── README.md                   # This file
```

---

## 🏛️ Domain-Driven Design Architecture (order-service)

The `order-service` is fully refactored using **DDD** with a **Ports and Adapters** (Hexagonal) pattern. Think of it in layers:

```
┌─────────────────────────────────────────────────────┐
│                  interface/grpc/                    │  gRPC handler (thin adapter)
│   Converts proto ↔ DTO, delegates to application   │
├─────────────────────────────────────────────────────┤
│                  application/                       │  Use cases (commands + queries)
│   Orchestrates domain objects, calls port interfaces│
├─────────────────────────────────────────────────────┤
│                    domain/                          │  Pure business logic
│   Aggregates, value objects, events, repo interfaces│
├─────────────────────────────────────────────────────┤
│                 infrastructure/                     │  Technical implementations
│   Postgres, RabbitMQ, gRPC client (implements ports)│
└─────────────────────────────────────────────────────┘
```

**Key DDD Concepts Applied:**

| Concept        | Implementation                                      | Purpose                                     |
| -------------- | --------------------------------------------------- | ------------------------------------------- |
| Aggregate Root | `Order` in `domain/aggregate/order.go`              | Single entry point for all order mutations  |
| Value Objects  | `Money`, `OrderID`, `OrderStatus`, `SKU`            | Immutable, self-validating domain concepts  |
| Domain Events  | `OrderCreatedEvent`, `OrderStatusUpdatedEvent`      | Decouple side effects from business logic   |
| Repository     | `OrderRepository` interface in `domain/repository/` | Abstracts persistence from domain           |
| Ports          | Interfaces in `application/port/`                   | Define what the app needs without how       |
| Adapters       | Everything in `infrastructure/`                     | Concrete implementations of port interfaces |

---

## ✨ Key Features

### 1. **Microservices Architecture**

- Three independent services with separate responsibilities
- Loose coupling through message queues
- Scalable and maintainable design

### 2. **gRPC Communication**

- High-performance binary protocol (Protocol Buffers)
- Type-safe API contracts
- Efficient inter-service calls (Order → Stock)

### 3. **Event-Driven Architecture**

- Asynchronous event publishing via RabbitMQ
- Decoupled services react to events independently
- Stock service and Analytics service consume order events

### 4. **Real-Time Analytics**

- Elasticsearch for fast search and aggregations
- Full-text search on customer names, SKUs, status
- Real-time statistics and reporting

### 5. **Database Management**

- PostgreSQL for reliable transactional data
- Proper schema design with foreign keys
- ACID compliance for critical operations

### 6. **Automated Stock Management**

- Automatic stock reservation on order creation
- Stock release on order cancellation
- Real-time inventory tracking

---

## 🚀 Services Overview

### 1. **Order Service (gRPC)**

- **Port:** 50061 (Docker) / 50051 (Manual)
- **Database:** PostgreSQL (orders_db)
- **Features:** Create orders, retrieve orders, update status, customer order history
- **Events:** Publishes order events to RabbitMQ

### 2. Stock Service (gRPC)

- **Port:** 50062 (Docker) / 50052 (Manual)
- **Database:** PostgreSQL (orders_db, stocks table)
- **Features:** Stock reservation, release, and inventory queries
- **Events:** Consumes order events to manage stock

### 3. Analytics Service (HTTP REST)

- **Port:** 8081 (Docker) / 8080 (Manual)
- **Database:** Elasticsearch
- **Features:** Search orders, aggregate statistics, health monitoring
- **Events:** Consumes order events to index in Elasticsearch

## Prerequisites

### Option 1: Docker (Recommended ⭐)

- Docker (with Docker Compose V2 plugin)
  - Use `docker compose` command (not `docker-compose`)
- `grpcurl` (for testing gRPC endpoints)
- `jq` (for JSON formatting)

### Option 2: Manual Setup

- Go 1.21+
- PostgreSQL
- RabbitMQ
- Elasticsearch
- `grpcurl` (for testing gRPC endpoints)
- `jq` (for JSON formatting)

---

## 🐳 Quick Start with Docker (Recommended)

### Simple 3-Step Setup

**Step 1: Clone and Navigate**

```bash
cd order-tracking-system
```

**Step 2: Start Everything**

```bash
docker compose up --build
```

This single command will:

- ✅ Pull and start PostgreSQL
- ✅ Pull and start RabbitMQ
- ✅ Pull and start Elasticsearch
- ✅ Build and start Order Service
- ✅ Build and start Stock Service
- ✅ Build and start Analytics Service
- ✅ Create database schema automatically
- ✅ Set up all networking

**Step 3: Verify Services**

```bash
# Check all containers are running
docker compose ps

# Expected output:
# order-tracking-postgres            running   5433/tcp
# order-tracking-rabbitmq            running   5673/tcp, 15673/tcp
# order-tracking-elasticsearch       running   9201/tcp, 9301/tcp
# order-tracking-order-service       running   50061/tcp
# order-tracking-stock-service       running   50062/tcp
# order-tracking-analytics-service   running   8081/tcp
```

### Access Services

**Docker Setup Ports (mapped to avoid conflicts with system services):**

- **Order Service (gRPC):** `localhost:50061`
- **Stock Service (gRPC):** `localhost:50062`
- **Analytics Service (HTTP):** `http://localhost:8081`
- **RabbitMQ Management UI:** `http://localhost:15673` (admin/admin)
- **Elasticsearch:** `http://localhost:9201`
- **PostgreSQL:** `localhost:5433` (admin/admin)

**Manual Setup Ports (default ports when running without Docker):**

- **Order Service (gRPC):** `localhost:50051`
- **Stock Service (gRPC):** `localhost:50052`
- **Analytics Service (HTTP):** `http://localhost:8080`
- **RabbitMQ Management UI:** `http://localhost:15672` (admin/admin)
- **Elasticsearch:** `http://localhost:9200`
- **PostgreSQL:** `localhost:5432` (admin/admin)

### Insert Sample Data

```bash
# Connect to PostgreSQL container
docker exec -it order-tracking-postgres psql -U admin -d orders_db

# Insert sample data
INSERT INTO customers (id, name) VALUES
('CUST-001', 'Md. Rakibul Kabir'),
('CUST-002', 'Arka Das'),
('CUST-003', 'Morshed Alam'),
('CUST-004', 'Emran Ahmed Emon');

INSERT INTO stocks (sku, quantity) VALUES
('LAPTOP-001', 50),
('MOUSE-001', 100),
('KEYBOARD-001', 75);

-- Exit with \q
```

### Useful Docker Commands

```bash
# Stop all services
docker compose down

# Stop and remove volumes (clean slate)
docker compose down -v

# View logs of all services
docker compose logs -f

# View logs of specific service
docker compose logs -f order-service

# Rebuild specific service
docker compose up --build order-service

# Restart specific service
docker compose restart order-service
```

### 🎮 Docker Manager Script (Interactive)

For an easier experience, use the interactive Docker manager:

```bash
./docker-manager.sh
```

This provides a menu-driven interface with options to:

- 🚀 Start/Stop/Restart services
- 🏗️ Rebuild services
- 📊 View status and logs
- 🧪 Insert sample data automatically
- 📦 Access PostgreSQL shell
- 🐰 Open RabbitMQ Management UI
- 🔍 Query Elasticsearch

---

## 🛠️ Manual Quick Start (Alternative)

> **Note:** This section is for running services **without Docker** directly on your host machine. Services will use their **default ports** (50051, 50052, 8080, etc.). This is different from Docker setup which uses mapped ports to avoid conflicts with system services.

### 1. Start Infrastructure Services

```bash
# Start PostgreSQL (default port 5432)
# Start RabbitMQ (default port 5672, management UI: 15672)
# Start Elasticsearch (default port 9200)
```

### 2. Initialize Database

```bash
# Connect to PostgreSQL
psql -h localhost -U admin -d orders_db

# Run schema files
\i order-service/db/schema.sql
\i stock-service/db/schema.sql

# Insert sample data
INSERT INTO customers (id, name) VALUES
('CUST-001', 'Md. Rakibul Kabir'),
('CUST-002', 'Arka Das'),
('CUST-003', 'Morshed Alam'),
('CUST-004', 'Emran Ahmed Emon');

INSERT INTO stocks (sku, quantity) VALUES
('LAPTOP-001', 50),
('MOUSE-001', 100),
('KEYBOARD-001', 75);
```

### 3. Start Services

Open three separate terminals:

**Terminal 1: Order Service**

```bash
cd order-service
go run main.go
```

**Terminal 2: Stock Service**

```bash
cd stock-service
go run main.go
```

**Terminal 3: Analytics Service**

```bash
cd analytics-service
go run main.go
```

### 4. Verify Services

```bash
# Check if all services are running
nc -z localhost 50051 && echo "✓ Order Service" || echo "✗ Order Service"
nc -z localhost 50052 && echo "✓ Stock Service" || echo "✗ Stock Service"
nc -z localhost 8080 && echo "✓ Analytics Service" || echo "✗ Analytics Service"
```

## Testing All Endpoints (11 Total)

> **Note:** The examples below use **Docker ports** (50061, 50062, 8081). If you're running services manually without Docker, replace with default ports (50051, 50052, 8080).

---

## ORDER SERVICE (gRPC) - 4 Methods

### 1. CreateOrder

Creates a new order with multiple items.

**Command:**

```bash
grpcurl -plaintext -d '{
  "customer_id": "CUST-001",
  "items": [
    {"sku": "LAPTOP-001", "quantity": 2, "unit_price": 999.99},
    {"sku": "MOUSE-001", "quantity": 1, "unit_price": 29.99}
  ]
}' localhost:50061 order.OrderService/CreateOrder
```

**Sample Response:**

```json
{
  "orderId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "pending",
  "message": "Order created successfully"
}
```

---

### 2. GetOrder

Retrieves order details by order ID.

**Command:**

```bash
grpcurl -plaintext -d '{
  "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}' localhost:50061 order.OrderService/GetOrder
```

**Sample Response:**

```json
{
  "orderId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "customerId": "CUST-001",
  "status": "pending",
  "totalAmount": 2029.97,
  "createdAt": "2026-03-02T10:30:00Z",
  "items": [
    { "sku": "LAPTOP-001", "quantity": 2, "unitPrice": 999.99 },
    { "sku": "MOUSE-001", "quantity": 1, "unitPrice": 29.99 }
  ]
}
```

---

### 3. UpdateOrderStatus

Updates the status of an existing order.

**Command:**

```bash
grpcurl -plaintext -d '{
  "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "shipped"
}' localhost:50061 order.OrderService/UpdateOrderStatus
```

**Sample Response:**

```json
{
  "orderId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "customerId": "CUST-001",
  "status": "shipped",
  "totalAmount": 2029.97,
  "createdAt": "2026-03-02T10:30:00Z",
  "items": [...]
}
```

**Valid Status Values:** `pending`, `confirmed`, `shipped`, `delivered`, `cancelled`

---

### 4. GetOrdersByCustomer

Retrieves all orders for a specific customer.

**Command:**

```bash
grpcurl -plaintext -d '{
  "customer_id": "CUST-001"
}' localhost:50061 order.OrderService/GetOrdersByCustomer
```

**Sample Response:**

```json
{
  "orders": [
    {
      "orderId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "customerId": "CUST-001",
      "status": "shipped",
      "totalAmount": 2029.97,
      "createdAt": "2026-03-02T10:30:00Z",
      "items": [...]
    },
    {
      "orderId": "b2c3d4e5-f6a7-8901-bcde-ef2345678901",
      "customerId": "CUST-001",
      "status": "delivered",
      "totalAmount": 999.99,
      "createdAt": "2026-03-01T14:20:00Z",
      "items": [...]
    }
  ]
}
```

---

## STOCK SERVICE (gRPC) - 3 Methods

### 5. GetStock

Retrieves current stock quantity for a SKU.

**Command:**

```bash
grpcurl -plaintext -d '{
  "sku": "LAPTOP-001"
}' localhost:50062 stock.StockService/GetStock
```

**Sample Response:**

```json
{
  "sku": "LAPTOP-001",
  "quantity": 48
}
```

---

### 6. ReserveStock

Reserves stock for an order (decrements quantity).

**Command:**

```bash
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 5
}' localhost:50062 stock.StockService/ReserveStock
```

**Sample Response:**

```json
{
  "success": true,
  "message": "Stock reserved successfully"
}
```

**Error Response (Insufficient Stock):**

```json
{
  "success": false,
  "message": "insufficient stock"
}
```

---

### 7. ReleaseStock

Releases reserved stock (increments quantity back).

**Command:**

```bash
grpcurl -plaintext -d '{
  "sku": "MOUSE-001",
  "quantity": 3
}' localhost:50062 stock.StockService/ReleaseStock
```

**Sample Response:**

```json
{
  "success": true,
  "message": "Stock released successfully"
}
```

---

## ANALYTICS SERVICE (HTTP REST) - 4 Endpoints

### 8. Search Orders (by customer name)

Search orders by customer name, SKU, or status.

**Command:**

```bash
curl -s "http://localhost:8081/search?customer=Rakib" | jq
```

**Sample Response:**

```json
{
  "total": 2,
  "orders": [
    {
      "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "customer_id": "CUST-001",
      "customer_name": "Md. Rakibul Kabir",
      "status": "shipped",
      "total_amount": 2029.97,
      "created_at": "2026-03-02T10:30:00Z",
      "items": [
        { "sku": "LAPTOP-001", "quantity": 2 },
        { "sku": "MOUSE-001", "quantity": 1 }
      ]
    }
  ]
}
```

---

### 9. Search Orders (by SKU)

**Command:**

```bash
curl -s "http://localhost:8081/search?sku=LAPTOP-001" | jq
```

**Sample Response:**

```json
{
  "total": 3,
  "orders": [
    {
      "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "customer_name": "Md. Rakibul Kabir",
      "items": [{ "sku": "LAPTOP-001", "quantity": 2 }]
    }
  ]
}
```

---

### 10. Search Orders (by status)

**Command:**

```bash
curl -s "http://localhost:8081/search?status=shipped" | jq
```

---

### 11. Aggregate Orders by Status

Get count of orders grouped by status.

**Command:**

```bash
curl -s "http://localhost:8081/aggregate/status" | jq
```

**Sample Response:**

```json
{
  "aggregations": [
    { "status": "pending", "count": 2 },
    { "status": "shipped", "count": 3 },
    { "status": "delivered", "count": 5 },
    { "status": "confirmed", "count": 1 }
  ]
}
```

---

### 12. Aggregate Orders by Customer

Get count of orders grouped by customer.

**Command:**

```bash
curl -s "http://localhost:8081/aggregate/customer" | jq
```

**Sample Response:**

```json
{
  "aggregations": [
    {
      "customer_id": "CUST-001",
      "customer_name": "Md. Rakibul Kabir",
      "count": 5
    },
    { "customer_id": "CUST-002", "customer_name": "Arka Das", "count": 3 },
    { "customer_id": "CUST-003", "customer_name": "Morshed Alam", "count": 2 }
  ]
}
```

---

### Bonus: Health Check

**Command:**

```bash
curl -s "http://localhost:8081/health" | jq
```

**Sample Response:**

```json
{
  "status": "ok"
}
```

---

## Complete Test Flow

Here's a complete workflow to test the entire system:

```bash
# 1. Create a new order
ORDER_ID=$(grpcurl -plaintext -d '{"customer_id":"CUST-001","items":[{"sku":"LAPTOP-001","quantity":1,"unit_price":999.99}]}' localhost:50061 order.OrderService/CreateOrder | grep -o '"orderId":"[^"]*"' | cut -d'"' -f4)

# 2. Get the order details
grpcurl -plaintext -d "{\"order_id\":\"$ORDER_ID\"}" localhost:50061 order.OrderService/GetOrder

# 3. Check stock was reserved
grpcurl -plaintext -d '{"sku":"LAPTOP-001"}' localhost:50062 stock.StockService/GetStock

# 4. Update order status
grpcurl -plaintext -d "{\"order_id\":\"$ORDER_ID\",\"status\":\"shipped\"}" localhost:50061 order.OrderService/UpdateOrderStatus

# 5. Wait 2 seconds for Elasticsearch indexing
sleep 2

# 6. Search in analytics
curl -s "http://localhost:8081/search?customer=Rakib" | jq

# 7. Check aggregations
curl -s "http://localhost:8081/aggregate/status" | jq
```

---

## Database Queries

View data directly from PostgreSQL:

**Docker Setup:**

```bash
# Use docker exec to connect
docker exec -it order-tracking-postgres psql -U admin -d orders_db

# View all customers
SELECT * FROM customers;

# View all orders
SELECT id, customer_id, status, total_amount FROM orders;

# View order items
SELECT * FROM order_items;

# View stock levels
SELECT sku, quantity FROM stocks;
```

**Manual Setup:**

```bash
# View all customers
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT * FROM customers;"

# View all orders
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT id, customer_id, status, total_amount FROM orders;"

# View order items
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT * FROM order_items;"

# View stock levels
PGPASSWORD=admin psql -h localhost -U admin -d orders_db -c "SELECT sku, quantity FROM stocks;"
```

---

## Elasticsearch Queries

**Docker Setup:**

```bash
# Get count of indexed orders
curl -s "http://localhost:9201/orders/_count" | jq

# View all indexed orders
curl -s "http://localhost:9201/orders/_search?size=10" | jq
```

**Manual Setup:**

```bash
# Get count of indexed orders
curl -s "http://localhost:9200/orders/_count" | jq

# View all indexed orders
curl -s "http://localhost:9200/orders/_search?size=10" | jq
```

---

## Environment Configuration

### Order Service (.env)

```env
DB_URL=postgres://admin:admin@localhost:5432/orders_db
RABBITMQ_URL=amqp://admin:admin@localhost:5672/
GRPC_PORT=50051
STOCK_SERVICE_ADDR=localhost:50052
```

### Stock Service (.env)

```env
DB_URL=postgres://admin:admin@localhost:5432/orders_db
RABBITMQ_URL=amqp://admin:admin@localhost:5672/
GRPC_PORT=50052
```

### Analytics Service (.env)

```env
ELASTICSEARCH_URL=http://localhost:9200
RABBITMQ_URL=amqp://admin:admin@localhost:5672/
HTTP_PORT=8080
```

---

## Troubleshooting

### Services not starting?

```bash
# Check if ports are already in use
lsof -i :50051  # Order Service
lsof -i :50052  # Stock Service
lsof -i :8080   # Analytics Service
```

### Database connection issues?

```bash
# Test PostgreSQL connection
psql -h localhost -U admin -d orders_db -c "SELECT 1;"
```

### RabbitMQ issues?

```bash
# Check RabbitMQ status
rabbitmqctl status

# View queues
rabbitmqctl list_queues
```

### Elasticsearch issues?

```bash
# Check Elasticsearch health
curl -s "http://localhost:9200/_cluster/health" | jq
```

---

## Technologies Stack

| Technology       | Version | Purpose                                      |
| ---------------- | ------- | -------------------------------------------- |
| Go               | 1.21+   | Primary programming language                 |
| gRPC             | -       | Inter-service communication protocol         |
| Protocol Buffers | 3       | API schema and serialization                 |
| PostgreSQL       | 14+     | Relational database for transactional data   |
| RabbitMQ         | 3.x     | Message broker for event-driven architecture |
| Elasticsearch    | 8.x     | Search engine and analytics platform         |

---

## Learning Outcomes

This project demonstrates practical knowledge of:

- Building microservices with independent deployment
- Implementing gRPC for efficient service-to-service communication
- Designing event-driven systems with message queues
- Integrating Elasticsearch for advanced search capabilities
- Managing distributed transactions and data consistency
- Working with Protocol Buffers for API contracts
- Building RESTful APIs alongside gRPC services

---

## Summary

**Total Operations: 11**

- Order Service: 4 gRPC methods
- Stock Service: 3 gRPC methods
- Analytics Service: 4 HTTP endpoints

All services work together to provide a complete order tracking solution with real-time analytics and automated stock management.
