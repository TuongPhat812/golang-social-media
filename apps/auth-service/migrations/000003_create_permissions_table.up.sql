-- Migration: Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    resource TEXT NOT NULL,
    action TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);

-- Updated_at trigger
CREATE TRIGGER update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

