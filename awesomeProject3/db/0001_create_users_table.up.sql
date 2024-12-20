CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100),
                       email VARCHAR(100) UNIQUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
