CREATE TABLE notification_subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

CREATE INDEX idx_notification_subscriptions_user_id ON notification_subscriptions(user_id);
CREATE INDEX idx_notification_subscriptions_enabled ON notification_subscriptions(enabled);