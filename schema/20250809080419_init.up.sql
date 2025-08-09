CREATE TABLE orders
(
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL UNIQUE,
    entry VARCHAR(255) NOT NULL,
    delivery_uid VARCHAR(255) REFERENCES deliveries(uid),
    --payment_transaction VARCHAR(255) REFERENCES payments(transaction),
    items --?,
    locale VARCHAR(255) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service TEXT NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP DEFAULT NOW(),
    oof_shard VARCHAR(255) NOT NULL
);

CREATE TABLE payments
(
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    transaction VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    payment_dt INT NOT NULL,
    bank VARCHAR (255) NOT NULL,
    delivery_cost INT DEFAULT 0,
    goods_total INT DEFAULT 1,
    custom_fee INT DEFAULT 0
);

CREATE TABLE items
(
    chrt_id INT PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INT DEFAULT 0,
    size VARCHAR(255) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INT DEFAULT 202
)