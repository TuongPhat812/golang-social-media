-- Rollback: Drop permissions table
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP TABLE IF EXISTS permissions;

