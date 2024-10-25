package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"

	"github.com/google/uuid"
)

type SupplierRepository interface {
	Create(ctx context.Context, req requests.CreateSupplierRequest) (*models.Supplier, error)
	Update(ctx context.Context, id uuid.UUID, req requests.UpdateSupplierRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error)
	List(ctx context.Context, limit, offset int) ([]models.Supplier, int64, error)
	GetByEmail(ctx context.Context, email string) (*models.Supplier, error)
}
