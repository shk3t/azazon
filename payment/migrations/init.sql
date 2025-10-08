CREATE TABLE processed_payment (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL UNIQUE
);

CREATE TABLE order_confirmed_outbox (
    id SERIAL PRIMARY KEY,
    processed BOOL NOT NULL,
    msg BYTEA NOT NULL
);

CREATE TABLE order_canceled_outbox (
    id SERIAL PRIMARY KEY,
    processed BOOL NOT NULL,
    msg BYTEA NOT NULL
);
