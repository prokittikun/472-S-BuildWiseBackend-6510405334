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
	"strings"

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

// List retrieves all jobs
func (r *jobRepository) List(ctx context.Context) (*responses.JobListResponse, error) {
	var jobs []models.Job
	query := `SELECT * FROM Job`
	err := r.db.SelectContext(ctx, &jobs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	var jobList []responses.JobResponse
	for _, job := range jobs {
		jobList = append(jobList, responses.JobResponse{
			JobID:       job.JobID,
			Name:        job.Name,
			Description: job.Description.String,
			Unit:        job.Unit,
		})
	}

	return &responses.JobListResponse{Jobs: jobList}, nil
}

// CreateJob creates a new job without materials
func (r *jobRepository) Create(ctx context.Context, req requests.CreateJobRequest) (*responses.JobResponse, error) {
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
		return &responses.JobResponse{
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

func (r *jobRepository) AddJobMaterial(ctx context.Context, jobID uuid.UUID, req requests.AddJobMaterialRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertJobMaterialQuery := `
		INSERT INTO Job_material (
			job_id, material_id, quantity
		) VALUES (
			:job_id, :material_id, :quantity
		) ON CONFLICT (job_id, material_id) 
		DO UPDATE SET quantity = Job_material.quantity + EXCLUDED.quantity`

	getProjectsQuery := `
		SELECT DISTINCT b.boq_id, b.status
		FROM boq_job bj 
		JOIN boq b ON b.boq_id = bj.boq_id 
		JOIN project p ON p.project_id = b.project_id 
		WHERE bj.job_id = $1`

	type BOQInfo struct {
		BOQID  uuid.UUID `db:"boq_id"`
		Status string    `db:"status"`
	}
	var boqs []BOQInfo

	err = tx.SelectContext(ctx, &boqs, getProjectsQuery, jobID)
	if err != nil {
		return fmt.Errorf("failed to get associated projects: %w", err)
	}

	// Insert job materials
	for _, material := range req.Materials {
		params := map[string]interface{}{
			"job_id":      jobID,
			"material_id": material.MaterialID,
			"quantity":    material.Quantity,
		}

		fmt.Print(params)

		_, err := tx.NamedExecContext(ctx, insertJobMaterialQuery, params)
		fmt.Print(err)

		if err != nil {

			return fmt.Errorf("failed to add material: %w", err)
		}

		// 14.4 For each draft BOQ, create material price log entries
		for _, boq := range boqs {
			if boq.Status == "draft" {
				insertPriceLogQuery := `
					INSERT INTO Material_price_log (
						material_id, boq_id, supplier_id, actual_price, estimated_price, job_id, quantity, updated_at
					) VALUES (
						$1, $2, NULL, NULL, NULL, $3, $4, CURRENT_TIMESTAMP
					)`

				_, err = tx.ExecContext(ctx, insertPriceLogQuery, material.MaterialID, boq.BOQID, jobID, material.Quantity)
				if err != nil {
					return fmt.Errorf("failed to create material price log: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *jobRepository) DeleteJobMaterial(ctx context.Context, jobID uuid.UUID, materialID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	checkUsageQuery := `
		SELECT DISTINCT p.name as project_name, b.boq_id, b.status
		FROM job_material jm 
		JOIN boq_job bj ON bj.job_id = jm.job_id 
		JOIN boq b ON b.boq_id = bj.boq_id 
		JOIN project p ON p.project_id = b.project_id 
		WHERE jm.material_id = $1`

	type ProjectUsage struct {
		ProjectName string    `db:"project_name"`
		BOQID       uuid.UUID `db:"boq_id"`
		Status      string    `db:"status"`
	}
	var usages []ProjectUsage

	err = tx.SelectContext(ctx, &usages, checkUsageQuery, materialID)
	if err != nil {
		return fmt.Errorf("failed to check material usage: %w", err)
	}
	fmt.Print(usages)
	if len(usages) > 0 {
		var projectNames []string
		for _, usage := range usages {
			projectNames = append(projectNames, usage.ProjectName)
		}
		return fmt.Errorf("material is used in following projects: %s", strings.Join(projectNames, ", "))
	}

	deleteJobMaterialQuery := `
		DELETE FROM Job_material 
		WHERE job_id = $1 AND material_id = $2`

	result, err := tx.ExecContext(ctx, deleteJobMaterialQuery, jobID, materialID)
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

	deletePriceLogsQuery := `
		DELETE FROM material_price_log 
		WHERE job_id = $1 
		AND material_id = $2 
		AND boq_id IN (
			SELECT b.boq_id 
			FROM boq b 
			JOIN boq_job bj ON b.boq_id = bj.boq_id 
			WHERE bj.job_id = $1 
			AND b.status = 'draft'
		)`

	_, err = tx.ExecContext(ctx, deletePriceLogsQuery, jobID, materialID)
	if err != nil {
		return fmt.Errorf("failed to delete material price logs: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

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
