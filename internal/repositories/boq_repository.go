package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type BOQRepository interface {
	Create(ctx context.Context, req requests.CreateBOQRequest) (*models.BOQ, error)             //don't use
	GetByID(ctx context.Context, id uuid.UUID) (*models.BOQ, error)                             //don't use
	GetByIDWithProject(ctx context.Context, id uuid.UUID) (*models.BOQ, *models.Project, error) //don't use
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.BOQ, error)               //don't use
	GetBoqWithProject(ctx context.Context, projectID uuid.UUID) (*responses.BOQResponse, error)
}
