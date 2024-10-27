package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type MaterialRepository interface {
	Create(ctx context.Context, req requests.CreateMaterialRequest) (*models.Material, error)
	Update(ctx context.Context, materialID string, req requests.UpdateMaterialRequest) error
	Delete(ctx context.Context, materialID string) error
	GetByID(ctx context.Context, materialID string) (*models.Material, error)
	List(ctx context.Context) ([]models.Material, error)

	GetMaterialPricesByProjectID(ctx context.Context, projectID uuid.UUID) ([]MaterialPriceInfo, error)
	UpdateEstimatedPrices(ctx context.Context, boqID uuid.UUID, materialID string, estimatedPrice float64) error
	GetBOQStatus(ctx context.Context, boqID uuid.UUID) (string, error)
}

type MaterialPriceInfo struct {
	MaterialID     string          `db:"material_id"`
	Name           string          `db:"name"`
	TotalQuantity  float64         `db:"qty_all_material_in_all_job"`
	Unit           string          `db:"unit"`
	EstimatedPrice sql.NullFloat64 `db:"estimated_price"`
	AvgActualPrice sql.NullFloat64 `db:"avg_actual_price"`
	ActualPrice    sql.NullFloat64 `db:"actual_price"`
	SupplierName   sql.NullString  `db:"supplier_name"`
}
