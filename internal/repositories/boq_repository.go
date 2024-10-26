package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"

	"github.com/google/uuid"
)

type BOQRepository interface {
	Create(ctx context.Context, req requests.CreateBOQRequest) (*models.BOQ, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BOQ, error)
	GetByIDWithProject(ctx context.Context, id uuid.UUID) (*models.BOQ, *models.Project, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.BOQ, error)
}
