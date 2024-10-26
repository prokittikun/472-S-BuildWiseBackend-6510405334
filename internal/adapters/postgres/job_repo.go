// repository/job_repository.go
package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type jobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) repositories.JobRepository {
	return &jobRepository{db: db}
}

// GetByID retrieves a job by its ID
func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	var job models.Job
	query := `SELECT * FROM Job WHERE job_id = $1`
	err := r.db.GetContext(ctx, &job, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("job not found")
		}
		return nil, fmt.Errorf("failed to get job by ID: %w", err)
	}
	return &job, nil
}

// CreateJob creates a new job without materials
func (r *jobRepository) Create(ctx context.Context, req requests.CreateJobRequest) (*responses.JobModelResponse, error) {
	job := &models.Job{
		JobID:       uuid.New(),
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Unit:        req.Unit,
	}

	query := `
        INSERT INTO Job (
            job_id, name, description, unit
        ) VALUES (
            :job_id, :name, :description, :unit
        ) RETURNING *`

	rows, err := r.db.NamedQueryContext(ctx, query, job)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(job)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		return &responses.JobModelResponse{
			JobID:       job.JobID,
			Name:        job.Name,
			Description: job.Description.String,
			Unit:        job.Unit,
		}, nil
	}

	return nil, errors.New("failed to create job: no rows returned")
}

func (r *jobRepository) Update(ctx context.Context, id uuid.UUID, req requests.UpdateJobRequest) error {
	query := `
        UPDATE Job SET 
            name = :name,
            description = :description,
            unit = :unit
        WHERE job_id = :job_id`

	params := map[string]interface{}{
		"job_id":      id,
		"name":        req.Name,
		"description": req.Description,
		"unit":        req.Unit,
	}

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("job not found")
	}

	return nil
}

// AddJobMaterial adds new material to a job
func (r *jobRepository) AddJobMaterial(ctx context.Context, jobID uuid.UUID, req requests.AddJobMaterialRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	fmt.Print("AddJobMaterial")

	query := `
		INSERT INTO Job_material (
			job_id, material_id, quantity
		) VALUES (
			:job_id, :material_id, :quantity
		) ON CONFLICT (job_id, material_id) 
		DO UPDATE SET quantity = Job_material.quantity + EXCLUDED.quantity`

	for _, material := range req.Materials {
		params := map[string]interface{}{
			"job_id":      jobID,
			"material_id": material.MaterialID,
			"quantity":    material.Quantity,
		}

		_, err := tx.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("failed to add material: %w", err)
		}
	}

	return tx.Commit()
}

func (r *jobRepository) DeleteJobMaterial(ctx context.Context, jobID uuid.UUID, materialID string) error {
	query := `
		DELETE FROM Job_material 
		WHERE job_id = $1 AND material_id = $2`

	result, err := r.db.ExecContext(ctx, query, jobID, materialID)
	if err != nil {
		return fmt.Errorf("failed to delete job material: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("job material not found")
	}

	return nil
}

// UpdateJobMaterialQuantity updates the quantity of a specific material in a job
func (r *jobRepository) UpdateJobMaterialQuantity(ctx context.Context, jobID uuid.UUID, req requests.UpdateJobMaterialQuantityRequest) error {
	query := `
		UPDATE Job_material 
		SET quantity = :quantity
		WHERE job_id = :job_id AND material_id = :material_id`

	params := map[string]interface{}{
		"job_id":      jobID,
		"material_id": req.MaterialID,
		"quantity":    req.Quantity,
	}

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update material quantity: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("job material not found")
	}

	return nil
}
