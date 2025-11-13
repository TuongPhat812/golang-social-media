CREATE KEYSPACE IF NOT EXISTS notification_service
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE notification_service;

CREATE TABLE IF NOT EXISTS notification_users (
    user_id text PRIMARY KEY,
    email text,
    name text,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS notifications_by_user (
    user_id text,
    created_at timestamp,
    notification_id uuid,
    type text,
    title text,
    body text,
    metadata map<text, text>,
    PRIMARY KEY (user_id, created_at, notification_id)
) WITH CLUSTERING ORDER BY (created_at DESC, notification_id DESC);

