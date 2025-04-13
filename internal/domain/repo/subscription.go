package repo

import (
	"context"
	"errors"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription models.Subscription) error
	GetAllSubscriptions(ctx context.Context, groupID int64) ([]*models.Subscription, error)
	DeleteSubscription(ctx context.Context, id string) error
	FindSubscription(ctx context.Context, groupID int64, email string) (*models.Subscription, error)
	FindOtherSubscription(ctx context.Context, groupID int64) (*models.Subscription, error)
}

var ErrSubscriptionNotFound = errors.New("subscription not found")
