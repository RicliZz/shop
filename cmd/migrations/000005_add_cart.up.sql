CREATE TABLE IF NOT EXISTS cart(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id),
    product_id int NOT NULL REFERENCES products(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);