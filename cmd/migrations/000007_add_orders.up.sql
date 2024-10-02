CREATE TYPE status AS ENUM('pending', 'on the way', 'delivered');

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    status status NOT NULL DEFAULT 'pending',
    total DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP default current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users (id)
)