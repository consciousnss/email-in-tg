package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID   int64  `bson:"_id" validate:"required"`
	Type string `bson:"type"`

	Title string `bson:"title"`

	Login *EmailLogin `bson:"login"`

	IsActive bool `bson:"is_active"`
}

type EmailLogin struct {
	Email    string `bson:"email" validate:"required,email"`
	Password string `bson:"password" validate:"required"`
}

type Subscription struct {
	ID          primitive.ObjectID `bson:"_id"`
	SenderEmail *string            `bson:"sender_email" validate:"omitempty,required_if=OtherSenders true,email"`
	GroupID     int64              `bson:"group_id" validate:"required"`
	ThreadID    int                `bson:"thread_id"`

	// if set to true, will be used if no SenderEmail matches found
	OtherSenders bool `bson:"other_senders"`
}
