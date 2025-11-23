-- Migration: Create users table for chat service
-- This table stores replicated user data from auth-service
-- Allows chat service to operate independently without querying auth-service

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on email for lookups (optional, if needed)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

