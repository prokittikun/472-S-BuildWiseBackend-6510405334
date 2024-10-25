package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/requests"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, req requests.CreateProjectRequest) (*models.Project, error)
	Update(ctx context.Context, id uuid.UUID, req requests.UpdateProjectRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetByIDWithClient(ctx context.Context, id uuid.UUID) (*models.Project, *models.Client, error)
	List(ctx context.Context) ([]models.Project, error)
}
