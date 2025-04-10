package repo

import (
	"context"

	"github.com/un1uckyyy/email-in-tg/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	db "github.com/un1uckyyy/email-in-tg/pkg/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	groupCollection         = "group"
	subscriptionsCollection = "subscriptions"
)

type Repo struct {
	db *mongo.Database
}

func NewRepo(db *db.Mongo) *Repo {
	return &Repo{db: db.DB}
}

func (r *Repo) CreateGroup(ctx context.Context, group models.Group) error {
	coll := r.db.Collection(groupCollection)

	_, err := coll.InsertOne(ctx, group)
	if mongo.IsDuplicateKeyError(err) {
		logger.Info("group already registered!")
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAllActiveGroups(ctx context.Context) ([]*models.Group, error) {
	coll := r.db.Collection(groupCollection)

	filter := bson.M{"is_active": true}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group
	if err = cur.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *Repo) SetEmailLogin(
	ctx context.Context,
	groupID int64,
	login models.EmailLogin,
) error {
	coll := r.db.Collection(groupCollection)

	filter := bson.M{"_id": groupID}
	update := bson.M{
		"$set": bson.M{
			"login": login,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) SetGroupActivity(ctx context.Context, groupID int64, activity bool) error {
	coll := r.db.Collection(groupCollection)
	filter := bson.M{"_id": groupID}
	update := bson.M{
		"$set": bson.M{
			"is_active": activity,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreateSubscription(ctx context.Context, subscription models.Subscription) error {
	coll := r.db.Collection(subscriptionsCollection)

	subscription.ID = primitive.NewObjectID()

	_, err := coll.InsertOne(ctx, subscription)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAllSubscriptions(ctx context.Context, groupID int64, threadID int) ([]*models.Subscription, error) {
	coll := r.db.Collection(subscriptionsCollection)

	filter := bson.M{"group_id": groupID, "thread_id": threadID}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var subscriptions []*models.Subscription
	if err = cur.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *Repo) DeleteSubscription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := r.db.Collection(subscriptionsCollection)

	filter := bson.M{"_id": objID}
	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) FindSubscription(ctx context.Context, groupID int64, email string) (*models.Subscription, error) {
	coll := r.db.Collection(subscriptionsCollection)

	filter := bson.M{"group_id": groupID, "sender_email": email}

	var subscription models.Subscription
	err := coll.FindOne(ctx, filter).Decode(&subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}
