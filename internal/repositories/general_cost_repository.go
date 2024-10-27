package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type GeneralCostRepository interface {
	Create(ctx context.Context, generalCost *models.GeneralCost) (*responses.GeneralCostResponse, error)
	GetByBOQID(ctx context.Context, boqID uuid.UUID) (*responses.GeneralCostListResponse, error)
	GetByID(ctx context.Context, gID uuid.UUID) (*models.GeneralCost, error)
	Update(ctx context.Context, gID uuid.UUID, req requests.UpdateGeneralCostRequest) error
	GetType(ctx context.Context) ([]models.Type, error)
}
