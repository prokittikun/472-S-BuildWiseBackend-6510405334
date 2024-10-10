package responses

import (
	"boonkosang/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type JobMaterialResponse struct {
	MaterialName  string  `json:"material_name"`
	Quantity      float64 `json:"quantity"`
	Type          string  `json:"type"`
	UnitOfMeasure string  `json:"unit_of_measure"`
}

type JobResponse struct {
	JobID       uuid.UUID             `json:"job_id"`
	Description string                `json:"description"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Materials   []JobMaterialResponse `json:"materials"`
}

func NewJobResponse(job *models.Job) JobResponse {
	materials := make([]JobMaterialResponse, len(job.Materials))
	for i, m := range job.Materials {
		materials[i] = JobMaterialResponse{
			MaterialName:  m.MaterialName,
			Quantity:      m.Quantity,
			Type:          m.Type,
			UnitOfMeasure: m.UnitOfMeasure,
		}
	}

	return JobResponse{
		JobID:       job.JobID,
		Description: job.Description,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
		Materials:   materials,
	}
}

func NewUpdateJobResponse(job *models.Job) JobResponse {
	return JobResponse{
		JobID:       job.JobID,
		Description: job.Description,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	}
}
