CREATE TABLE IF NOT EXISTS items
(
    uuid VARCHAR(100) PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL CHECK ( price >= 0 ),
    CONSTRAINT unique_product_name UNIQUE (product_name)

);

CREATE TABLE IF NOT EXISTS users
(
    uuid VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL ,
    password VARCHAR(100) NOT NULL ,
    coins INTEGER NOT NULL CHECK ( coins >= 0 ) DEFAULT 1000,
    CONSTRAINT unique_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY ,
    sender VARCHAR(100) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE ,
    recipient VARCHAR(100) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE ,
    amount INTEGER NOT NULL CHECK ( amount > 0 ),
    date_created TIMESTAMPTZ NOT NULL,
    CONSTRAINT sender_recipient_diff CHECK ( sender <> recipient )

);

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY ,
    user_uuid VARCHAR(100) REFERENCES users(uuid) ON DELETE CASCADE,
    items_uuid VARCHAR(100) REFERENCES items(uuid) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK ( quantity > 0 ),
    total_price INTEGER NOT NULL CHECK (total_price > 0),
    date_created TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS ownership (
    user_uuid VARCHAR(100) REFERENCES users(uuid) ON DELETE CASCADE,
    items_uuid VARCHAR(100) REFERENCES items(uuid) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK ( quantity > 0 )
)