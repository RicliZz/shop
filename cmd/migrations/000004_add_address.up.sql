CREATE TABLE IF NOT EXISTS address(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    city VARCHAR(30),
    street VARCHAR(30),
    house VARCHAR(10),
    apartment VARCHAR(10),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);