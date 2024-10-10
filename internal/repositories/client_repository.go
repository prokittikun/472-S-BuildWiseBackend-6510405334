package repositories

import (
	"boonkosang/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, client *models.Client) error
	ListClients(ctx context.Context) ([]*models.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	UpdateClient(ctx context.Context, client *models.Client) error
	DeleteClient(ctx context.Context, id uuid.UUID) error
}
