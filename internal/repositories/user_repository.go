package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user requests.RegisterRequest) error
}
