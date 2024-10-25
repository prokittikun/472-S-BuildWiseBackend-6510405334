package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"

	"github.com/google/uuid"
)

type ClientRepository interface {
	Create(ctx context.Context, req requests.CreateClientRequest) (*models.Client, error)
	Update(ctx context.Context, id uuid.UUID, req requests.UpdateClientRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	List(ctx context.Context, limit, offset int) ([]models.Client, int64, error)
	GetByEmail(ctx context.Context, email string) (*models.Client, error)
}
