package repo

import "errors"

var (
	ErrGroupNotFound        = errors.New("group not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
