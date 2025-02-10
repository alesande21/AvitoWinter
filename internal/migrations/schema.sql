CREATE TABLE IF NOT EXISTS items
(
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL CHECK ( price >= 0 ),
    CONSTRAINT unique_product_name UNIQUE (product_name)

);

CREATE TABLE IF NOT EXISTS users
(
    uuid VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL ,
    password VARCHAR(100) NOT NULL ,
    wallet INTEGER NOT NULL CHECK ( wallet >= 0 ),
    CONSTRAINT unique_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS users_transfer (
    id SERIAL PRIMARY KEY ,
    from_user VARCHAR(100) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE ,
    to_user VARCHAR(100) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE ,
    amount INTEGER NOT NULL CHECK ( amount > 0 ),
    date_created TIMESTAMPTZ NOT NULL
)

CREATE TABLE IF NOT EXISTS orders
(
    uuid VARCHAR(100) PRIMARY KEY,
    user_uuid VARCHAR(100) REFERENCES users(uuid) ON DELETE CASCADE,
    items_uuid VARCHAR(100) REFERENCES items(id) ON DELETE CASCADE,
    size INTEGER NOT NULL CHECK ( size > 0 ),
    total_price INTEGER NOT NULL CHECK (total_price > 0),
    date_created TIMESTAMPTZ NOT NULL
)