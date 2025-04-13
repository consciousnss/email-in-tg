package mongo

import (
	"context"
	"errors"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"github.com/un1uckyyy/email-in-tg/internal/domain/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	subscriptionsCollection = "subscriptions"
)

type subscriptionRepo struct {
	coll *mongo.Collection
}

var _ repo.SubscriptionRepository = (*subscriptionRepo)(nil)

func NewSubscriptionRepo(db *mongo.Database) repo.SubscriptionRepository {
	return &subscriptionRepo{coll: db.Collection(subscriptionsCollection)}
}

func (sr *subscriptionRepo) CreateSubscription(ctx context.Context, subscription models.Subscription) error {
	sub := toMongoSubscription(&subscription)

	_, err := sr.coll.InsertOne(ctx, sub)
	if err != nil {
		return err
	}
	return nil
}

func (sr *subscriptionRepo) GetAllSubscriptions(ctx context.Context, groupID int64) ([]*models.Subscription, error) {
	filter := bson.M{"group_id": groupID}
	cur, err := sr.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var subscriptions []*subscription
	if err = cur.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	domains := make([]*models.Subscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		domains = append(domains, fromMongoSubscription(s))
	}

	return domains, nil
}

func (sr *subscriptionRepo) DeleteSubscription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = sr.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (sr *subscriptionRepo) FindSubscription(
	ctx context.Context,
	groupID int64,
	email string,
) (*models.Subscription, error) {
	filter := bson.M{"group_id": groupID, "sender_email": email}

	var sub subscription
	err := sr.coll.FindOne(ctx, filter).Decode(&sub)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// trying to find subscription on other senders for that group
		s, err := sr.FindOtherSubscription(ctx, groupID)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repo.ErrSubscriptionNotFound
		}
		if err != nil {
			return nil, err
		}

		return s, nil
	}
	if err != nil {
		return nil, err
	}

	return fromMongoSubscription(&sub), nil
}

// FindOtherSubscription return subscription on other senders for group or error
func (sr *subscriptionRepo) FindOtherSubscription(ctx context.Context, groupID int64) (*models.Subscription, error) {
	filter := bson.M{"group_id": groupID, "other_senders": true}

	var sub subscription
	err := sr.coll.FindOne(ctx, filter).Decode(&sub)
	if err != nil {
		return nil, err
	}

	return fromMongoSubscription(&sub), nil
}
