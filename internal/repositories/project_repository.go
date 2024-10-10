package repositories

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.Project) error
	ListProjects(ctx context.Context) ([]*responses.ProjectResponse, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}
