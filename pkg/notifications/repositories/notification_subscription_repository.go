package repositories

import (
	"database/sql"
	"fmt"
	"tidebot/pkg/notifications/models"

	"github.com/labstack/echo/v4"
)

type NotificationSubscriptionRepository interface {
	CreateSubscription(userID int) error
	GetSubscriptionByUserID(userID int) (*models.NotificationSubscription, error)
	EnableSubscription(userID int) error
	DisableSubscription(userID int) error
	GetEnabledSubscriptions() ([]models.NotificationSubscription, error)
}

type notificationSubscriptionRepositoryImpl struct {
	db  *sql.DB
	log echo.Logger
}

func NewNotificationSubscriptionRepository(db *sql.DB, log echo.Logger) NotificationSubscriptionRepository {
	return &notificationSubscriptionRepositoryImpl{
		db:  db,
		log: log,
	}
}

func (r *notificationSubscriptionRepositoryImpl) CreateSubscription(userID int) error {
	query := `
		INSERT INTO notification_subscriptions (user_id, enabled) 
		VALUES (?, ?) 
		ON CONFLICT(user_id) DO UPDATE SET enabled = ?, updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := r.db.Exec(query, userID, true, true)
	if err != nil {
		r.log.Errorf("Failed to create subscription for user %d: %v", userID, err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	
	r.log.Infof("Created/enabled subscription for user %d", userID)
	return nil
}

func (r *notificationSubscriptionRepositoryImpl) GetSubscriptionByUserID(userID int) (*models.NotificationSubscription, error) {
	query := `
		SELECT id, user_id, enabled, created_at, updated_at 
		FROM notification_subscriptions 
		WHERE user_id = ?
	`
	
	var subscription models.NotificationSubscription
	err := r.db.QueryRow(query, userID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Enabled,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		r.log.Errorf("Failed to get subscription for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	
	return &subscription, nil
}

func (r *notificationSubscriptionRepositoryImpl) EnableSubscription(userID int) error {
	query := `
		UPDATE notification_subscriptions 
		SET enabled = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE user_id = ?
	`
	
	result, err := r.db.Exec(query, true, userID)
	if err != nil {
		r.log.Errorf("Failed to enable subscription for user %d: %v", userID, err)
		return fmt.Errorf("failed to enable subscription: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found for user %d", userID)
	}
	
	r.log.Infof("Enabled subscription for user %d", userID)
	return nil
}

func (r *notificationSubscriptionRepositoryImpl) DisableSubscription(userID int) error {
	query := `
		UPDATE notification_subscriptions 
		SET enabled = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE user_id = ?
	`
	
	result, err := r.db.Exec(query, false, userID)
	if err != nil {
		r.log.Errorf("Failed to disable subscription for user %d: %v", userID, err)
		return fmt.Errorf("failed to disable subscription: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found for user %d", userID)
	}
	
	r.log.Infof("Disabled subscription for user %d", userID)
	return nil
}

func (r *notificationSubscriptionRepositoryImpl) GetEnabledSubscriptions() ([]models.NotificationSubscription, error) {
	query := `
		SELECT id, user_id, enabled, created_at, updated_at 
		FROM notification_subscriptions 
		WHERE enabled = ?
	`
	
	rows, err := r.db.Query(query, true)
	if err != nil {
		r.log.Errorf("Failed to get enabled subscriptions: %v", err)
		return nil, fmt.Errorf("failed to get enabled subscriptions: %w", err)
	}
	defer rows.Close()
	
	var subscriptions []models.NotificationSubscription
	for rows.Next() {
		var subscription models.NotificationSubscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Enabled,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("Failed to scan subscription: %v", err)
			continue
		}
		subscriptions = append(subscriptions, subscription)
	}
	
	return subscriptions, nil
}