package mongo

import (
	"context"
	"errors"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"github.com/un1uckyyy/email-in-tg/internal/domain/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	groupCollection = "group"
)

type groupRepo struct {
	coll *mongo.Collection
}

var _ repo.GroupRepository = (*groupRepo)(nil)

func NewGroupRepo(db *mongo.Database) repo.GroupRepository {
	return &groupRepo{coll: db.Collection(groupCollection)}
}

func (gr *groupRepo) CreateGroup(ctx context.Context, group models.Group) error {
	g := toMongoGroup(&group)

	_, err := gr.coll.InsertOne(ctx, g)
	if mongo.IsDuplicateKeyError(err) {
		return repo.ErrGroupAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (gr *groupRepo) GetGroup(ctx context.Context, id int64) (*models.Group, error) {
	filter := bson.M{"_id": id}

	var g group
	err := gr.coll.FindOne(ctx, filter).Decode(&g)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, repo.ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}

	return fromMongoGroup(&g), nil
}

func (gr *groupRepo) GetAllActiveGroups(ctx context.Context) ([]*models.Group, error) {
	filter := bson.M{"is_active": true}
	cur, err := gr.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var groups []*group
	if err = cur.All(ctx, &groups); err != nil {
		return nil, err
	}

	domains := make([]*models.Group, 0, len(groups))
	for _, g := range groups {
		domains = append(domains, fromMongoGroup(g))
	}

	return domains, nil
}

func (gr *groupRepo) SetEmailLogin(
	ctx context.Context,
	groupID int64,
	login models.EmailLogin,
) error {
	filter := bson.M{"_id": groupID}
	update := bson.M{
		"$set": bson.M{
			"login": toMongoEmailLogin(&login),
		},
	}

	_, err := gr.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (gr *groupRepo) SetGroupActivity(ctx context.Context, groupID int64, activity bool) (*models.Group, error) {
	filter := bson.M{"_id": groupID}
	update := bson.M{
		"$set": bson.M{
			"is_active": activity,
		},
	}

	var g group
	err := gr.coll.FindOneAndUpdate(ctx, filter, update).Decode(&g)
	if err != nil {
		return nil, err
	}

	return fromMongoGroup(&g), nil
}
