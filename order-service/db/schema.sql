CREATE TABLE IF NOT EXISTS customers (
    id         VARCHAR(100) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id            VARCHAR(100) PRIMARY KEY,
    customer_id   VARCHAR(100) NOT NULL REFERENCES customers(id),
    total_amount  DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status        VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id          SERIAL PRIMARY KEY,
    order_id    VARCHAR(100) NOT NULL REFERENCES orders(id),
    sku         VARCHAR(100) NOT NULL,
    quantity    INTEGER NOT NULL,
    unit_price  DECIMAL(10, 2) NOT NULL
);