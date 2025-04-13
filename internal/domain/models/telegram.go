package models

type Group struct {
	ID   int64 `validate:"required"`
	Type string

	Title string

	Login *EmailLogin

	IsActive bool
}

type EmailLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type Subscription struct {
	ID          string
	SenderEmail *string `validate:"omitempty,required_if=OtherSenders false,email"`
	GroupID     int64   `validate:"required"`
	ThreadID    int

	// if set to true, will be used if no SenderEmail matches found
	OtherSenders bool
}
