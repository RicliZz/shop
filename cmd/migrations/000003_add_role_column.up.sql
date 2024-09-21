CREATE TYPE role AS ENUM('admin', 'user');
ALTER TABLE users ADD COLUMN role role DEFAULT 'user';