package mongo

import (
	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type group struct {
	ID   int64  `bson:"_id"`
	Type string `bson:"type"`

	Title string `bson:"title"`

	Login *emailLogin `bson:"login"`

	IsActive bool `bson:"is_active"`
}

type emailLogin struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type subscription struct {
	ID          primitive.ObjectID `bson:"_id"`
	SenderEmail *string            `bson:"sender_email"`
	GroupID     int64              `bson:"group_id"`
	ThreadID    int                `bson:"thread_id"`

	// if set to true, will be used if no SenderEmail matches found
	OtherSenders bool `bson:"other_senders"`
}

func toMongoGroup(g *models.Group) *group {
	return &group{
		ID:       g.ID,
		Type:     g.Type,
		Title:    g.Title,
		Login:    toMongoEmailLogin(g.Login),
		IsActive: g.IsActive,
	}
}

func fromMongoGroup(g *group) *models.Group {
	return &models.Group{
		ID:       g.ID,
		Type:     g.Type,
		Title:    g.Title,
		Login:    fromMongoEmailLogin(g.Login),
		IsActive: g.IsActive,
	}
}

func toMongoEmailLogin(login *models.EmailLogin) *emailLogin {
	return &emailLogin{
		Email:    login.Email,
		Password: login.Password,
	}
}

func fromMongoEmailLogin(login *emailLogin) *models.EmailLogin {
	return &models.EmailLogin{
		Email:    login.Email,
		Password: login.Password,
	}
}

func toMongoSubscription(s *models.Subscription) *subscription {
	return &subscription{
		ID:           primitive.NewObjectID(),
		SenderEmail:  s.SenderEmail,
		GroupID:      s.GroupID,
		ThreadID:     s.ThreadID,
		OtherSenders: s.OtherSenders,
	}
}

func fromMongoSubscription(s *subscription) *models.Subscription {
	return &models.Subscription{
		ID:           s.ID.Hex(),
		SenderEmail:  s.SenderEmail,
		GroupID:      s.GroupID,
		ThreadID:     s.ThreadID,
		OtherSenders: s.OtherSenders,
	}
}
