CREATE TABLE customers
(
    customer_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(127) NOT NULL,
    surname VARCHAR(127) NOT NULL,
    phone VARCHAR(15) NOT NULL UNIQUE,
    email VARCHAR(127) NOT NULL UNIQUE
);

CREATE TABLE deliveries
(
    id SERIAL PRIMARY KEY,
    region VARCHAR(15) NOT NULL,
    zip VARCHAR(15) NOT NULL,
    city VARCHAR(255) NOT NULL,
    street VARCHAR(255) NOT NULL,
    house VARCHAR (15) NOT NULL,
    flat VARCHAR (15)
);

CREATE TABLE payments
(
    payment_uid VARCHAR(255) PRIMARY KEY,
    transaction VARCHAR(31) NOT NULL UNIQUE,
    request_id VARCHAR(255),
    currency VARCHAR (3) DEFAULT 'RUB',
    provider VARCHAR (255) NOT NULL,
    amount INT DEFAULT 1,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INT DEFAULT 0,
    goods_total INT DEFAULT 1,
    custom_fee INT DEFAULT 0
);

CREATE TABLE items
(   
    item_id SERIAL PRIMARY KEY,
    nm_id INT,
    name VARCHAR(255),
    size VARCHAR(255) DEFAULT '0',
    brand VARCHAR(255) NOT NULL
);


CREATE TABLE orders
(
   order_uid VARCHAR(255) PRIMARY KEY,
   track_number VARCHAR(255) NOT NULL,
   entry_code VARCHAR(63) NOT NULL,
   internal_signature VARCHAR(255),
   shardkey VARCHAR(15),
   sm_id INT NOT NULL,
   date_created TIMESTAMP DEFAULT NOW(),
   oof_shard VARCHAR (15) NOT NULL,
   locale VARCHAR(3) NOT NULL,
   customer_id VARCHAR(255) REFERENCES customers(customer_uid),
   delivery_service VARCHAR(255),
   delivery_id INT REFERENCES deliveries(id),
   payment_id VARCHAR(255) REFERENCES payments(payment_uid)
);

CREATE TABLE order_items
(   
    item_id INT REFERENCES items(item_id) ON DELETE CASCADE,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INT,
    price INT NOT NULL,
    rid VARCHAR(63) NOT NULL UNIQUE,
    sale INT DEFAULT 0,
    total_price INT,
    status INT NOT NULL,
    PRIMARY KEY(item_id, order_uid)
);
