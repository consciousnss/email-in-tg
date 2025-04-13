package repo

import (
	"context"
	"errors"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group models.Group) error
	GetGroup(ctx context.Context, id int64) (*models.Group, error)
	GetAllActiveGroups(ctx context.Context) ([]*models.Group, error)
	SetEmailLogin(ctx context.Context, groupID int64, login models.EmailLogin) error
	SetGroupActivity(ctx context.Context, groupID int64, activity bool) (*models.Group, error)
}

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupAlreadyExists = errors.New("group already exists")
)
