CREATE TABLE "order" (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR(16) NOT NULL,
    address VARCHAR(64) NOT NULL,
    track VARCHAR(64) NOT NULL
);

CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES "order"(id) ON DELETE CASCADE
);
