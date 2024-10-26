package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/responses"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type boqRepository struct {
	db *sqlx.DB
}

func NewBOQRepository(db *sqlx.DB) repositories.BOQRepository {
	return &boqRepository{
		db: db,
	}
}

func (r *boqRepository) GetBoqWithProject(ctx context.Context, projectID uuid.UUID) (*responses.BOQResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var data struct {
		BoqID              uuid.UUID        `db:"boq_id"`
		ProjectID          uuid.UUID        `db:"project_id"`
		BOQStatus          models.BOQStatus `db:"boq_status"`
		SellingGeneralCost sql.NullFloat64  `db:"selling_general_cost"`
	}

	boqQuery := `
        SELECT 
            b.boq_id,
            b.project_id,
            b.status as boq_status,
            b.selling_general_cost
        FROM boq b
        JOIN project p ON p.project_id = b.project_id
        WHERE b.project_id = $1`

	err = tx.GetContext(ctx, &data, boqQuery, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new BOQ if it doesn't exist
			createBOQQuery := `
                INSERT INTO boq (project_id, status, selling_general_cost) 
                VALUES ($1, 'draft', NULL) 
                RETURNING boq_id, project_id, status as boq_status, selling_general_cost`

			err = tx.GetContext(ctx, &data, createBOQQuery, projectID)
			if err != nil {
				return nil, fmt.Errorf("failed to create BOQ: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to check BOQ existence: %w", err)
		}
	}

	// Convert to response struct
	response := &responses.BOQResponse{
		ID:                 data.BoqID,
		ProjectID:          data.ProjectID,
		Status:             data.BOQStatus,
		SellingGeneralCost: data.SellingGeneralCost.Float64,
	}

	fmt.Print(response)

	jobsQuery := `
   SELECT DISTINCT
	j.*
FROM job j
JOIN boq_job bj ON j.job_id = bj.job_id
WHERE bj.boq_id = $1
`

	var jobs []models.Job
	err = tx.SelectContext(ctx, &jobs, jobsQuery, data.BoqID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	var jobForResponse []responses.JobResponse
	for _, job := range jobs {
		jobForResponse = append(jobForResponse, responses.JobResponse{
			JobID:       job.JobID,
			Name:        job.Name,
			Description: job.Description.String,
			Unit:        job.Unit,
		})
	}

	response.Jobs = jobForResponse

	return response, nil
}
