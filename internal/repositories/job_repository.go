package repositories

import (
	"boonkosang/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type JobRepository interface {
	CreateJob(ctx context.Context, job *models.Job, materials []models.JobMaterial) error
	ListJobs(ctx context.Context) ([]*models.Job, error)
	GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error)
	UpdateJob(ctx context.Context, job *models.Job, materials []models.JobMaterial) error
	DeleteJob(ctx context.Context, id uuid.UUID) error
	GetJobMaterials(ctx context.Context, jobID uuid.UUID) ([]models.JobMaterial, error)
}
