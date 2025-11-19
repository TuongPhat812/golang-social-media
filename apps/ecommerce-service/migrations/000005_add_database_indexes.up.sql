-- Add indexes for frequently queried columns

-- Products table indexes
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_products_status_created_at ON products(status, created_at DESC);

-- Orders table indexes
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_user_id_created_at ON orders(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_user_id_status ON orders(user_id, status);

-- Order items table indexes
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id_product_id ON order_items(order_id, product_id);

-- Outbox table indexes (already exists, but ensure they're optimal)
-- These are already created in 000003 migration

-- Event store table indexes (already exists, but ensure they're optimal)
-- These are already created in 000004 migration
