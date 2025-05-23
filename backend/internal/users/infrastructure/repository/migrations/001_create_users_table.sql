-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,
    role VARCHAR(20) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    documents JSONB
);

-- Add indices for common queries
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
