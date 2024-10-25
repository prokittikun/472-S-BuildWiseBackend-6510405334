package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"
)

type MaterialRepository interface {
	Create(ctx context.Context, req requests.CreateMaterialRequest) (*models.Material, error)
	Update(ctx context.Context, materialID string, req requests.UpdateMaterialRequest) error
	Delete(ctx context.Context, materialID string) error
	GetByID(ctx context.Context, materialID string) (*models.Material, error)
	List(ctx context.Context) ([]models.Material, error)
}
