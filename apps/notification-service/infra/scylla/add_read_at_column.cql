-- Migration: Add read_at column to notifications_by_user table
-- This migration adds the read_at column to track when notifications are read
-- Note: If the column already exists, this will fail. That's expected behavior.

USE notification_service;

-- Add read_at column to existing table
-- If column already exists, you'll get an error - that's OK, just means migration was already applied
ALTER TABLE notifications_by_user ADD read_at timestamp;

