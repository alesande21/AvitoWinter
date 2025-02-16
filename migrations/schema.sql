CREATE TABLE IF NOT EXISTS items
(
    product_name VARCHAR(100) PRIMARY KEY,
    price INTEGER NOT NULL CHECK ( price >= 0 ),
    CONSTRAINT unique_product_name UNIQUE (product_name)
);

CREATE TABLE IF NOT EXISTS users
(
    username VARCHAR(100) PRIMARY KEY ,
    password VARCHAR(100) NOT NULL ,
    coins INTEGER NOT NULL CHECK ( coins >= 0 ) DEFAULT 1000,
    CONSTRAINT unique_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY ,
    sender VARCHAR(100) NOT NULL REFERENCES users(username) ON DELETE CASCADE ,
    recipient VARCHAR(100) NOT NULL REFERENCES users(username) ON DELETE CASCADE ,
    amount INTEGER NOT NULL CHECK ( amount > 0 ),
    date_created TIMESTAMPTZ NOT NULL,
    CONSTRAINT sender_recipient_diff CHECK ( sender <> recipient )

);

CREATE TABLE IF NOT EXISTS purchases
(
    id SERIAL PRIMARY KEY ,
    username VARCHAR(100) REFERENCES users(username) ON DELETE CASCADE,
    item VARCHAR(100) REFERENCES items(product_name) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK ( quantity > 0 ),
    total_price INTEGER NOT NULL CHECK (total_price > 0),
    date_created TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS ownership (
    username VARCHAR(100) REFERENCES users(username) ON DELETE CASCADE,
    item VARCHAR(100) REFERENCES items(product_name) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK ( quantity > 0 ),
    PRIMARY KEY (username, item)
);

INSERT INTO items (product_name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500)
ON CONFLICT (product_name) DO NOTHING;


