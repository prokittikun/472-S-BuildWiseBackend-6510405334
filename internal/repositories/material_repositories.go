package repositories

import (
	"boonkosang/internal/domain/models"
	"context"
)

type MaterialRepository interface {
	CreateMaterial(ctx context.Context, material *models.Material) error
	ListMaterials(ctx context.Context) ([]*models.Material, error)
	GetMaterialByName(ctx context.Context, name string) (*models.Material, error)
	UpdateMaterial(ctx context.Context, material *models.Material) error
	DeleteMaterial(ctx context.Context, name string) error
	GetMaterialPriceHistory(ctx context.Context, name string) ([]*models.MaterialPriceLog, error)
	MaterialExists(ctx context.Context, name string) (bool, error)
}
