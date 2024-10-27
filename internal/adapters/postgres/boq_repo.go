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

type boqRepository struct {
	db *sqlx.DB
}

func NewBOQRepository(db *sqlx.DB) repositories.BOQRepository {
	return &boqRepository{
		db: db,
	}
}

func (r *boqRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BOQ, error) {
	var boq models.BOQ
	query := `SELECT * FROM boq WHERE boq_id = $1`
	err := r.db.GetContext(ctx, &boq, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get BOQ: %w", err)
	}

	return &boq, nil
}

func (r *boqRepository) Approve(ctx context.Context, boqID uuid.UUID) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)

	}

	// Check BOQ status
	var status string
	checkStatusQuery := `SELECT status FROM boq WHERE boq_id = $1`
	err = tx.GetContext(ctx, &status, checkStatusQuery, boqID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("boq not found")
		}
		return fmt.Errorf("failed to get BOQ status: %w", err)
	}

	if status != "draft" {
		return errors.New("can only approve BOQ in draft status")

	}

	// Update BOQ status
	updateQuery := `UPDATE boq SET status = 'approved' WHERE boq_id = $1`
	_, err = tx.ExecContext(ctx, updateQuery, boqID)
	if err != nil {
		return fmt.Errorf("failed to update BOQ status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *boqRepository) GetBoqWithProject(ctx context.Context, projectID uuid.UUID) (*responses.BOQResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var data models.BOQ

	boqQuery := `
        SELECT  boq_id, project_id, status, selling_general_cost
		FROM Boq
		WHERE project_id = $1`

	err = tx.GetContext(ctx, &data, boqQuery, projectID)
	fmt.Print(data)
	fmt.Print(err)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Print("create new boq")
			// Create new BOQ if it doesn't exist
			createBOQQuery := `
                INSERT INTO Boq (project_id, status, selling_general_cost) 
                VALUES (:project_id, 'draft', NULL) 
                RETURNING boq_id, project_id, status, selling_general_cost`

			row, err := r.db.NamedQueryContext(ctx, createBOQQuery, map[string]interface{}{
				"project_id": projectID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create new BOQ: %w", err)
			}

			if row.Next() {
				err = row.StructScan(&data)
				if err != nil {
					return nil, fmt.Errorf("failed to scan BOQ: %w", err)
				}
			}

			if err := row.Close(); err != nil {
				return nil, fmt.Errorf("failed to close row: %w", err)
			}

		} else {
			return nil, fmt.Errorf("failed to check BOQ existence: %w", err)
		}
	}
	fmt.Print(data)

	// Convert to response struct
	response := &responses.BOQResponse{
		ID:                 data.BOQID, // Assuming the correct field name is BOQID
		ProjectID:          data.ProjectID,
		Status:             data.Status, // Assuming the correct field name is Status
		SellingGeneralCost: data.SellingGeneralCost.Float64,
	}

	jobsQuery := `
   SELECT DISTINCT
	j.*
FROM job j
JOIN boq_job bj ON j.job_id = bj.job_id
WHERE bj.boq_id = $1
`

	var jobs []models.Job
	err = tx.SelectContext(ctx, &jobs, jobsQuery, data.BOQID)
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

func (r *boqRepository) AddBOQJob(ctx context.Context, boqID uuid.UUID, req requests.BOQJobRequest) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check BOQ status
	var status string
	checkStatusQuery := `SELECT status FROM boq WHERE boq_id = $1`
	err = tx.GetContext(ctx, &status, checkStatusQuery, boqID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("boq not found")
		}
		return fmt.Errorf("failed to get BOQ status: %w", err)
	}

	if status != "draft" {
		return errors.New("can only add jobs to BOQ in draft status")
	}

	// 16.2 Insert into boq_job
	insertBOQJobQuery := `
        INSERT INTO boq_job (
            boq_id, job_id, quantity, labor_cost
        ) VALUES (
            $1, $2, $3, $4
        )`

	_, err = tx.ExecContext(ctx, insertBOQJobQuery,
		boqID,
		req.JobID,
		req.Quantity,
		req.LaborCost,
	)
	if err != nil {
		return fmt.Errorf("failed to add job to BOQ: %w", err)
	}

	// 16.3 Get materials for the job and add to material_price_log if not exists
	materialQuery := `
        SELECT material_id, quantity 
        FROM job_material 
        WHERE job_id = $1`

	type JobMaterial struct {
		MaterialID string  `db:"material_id"`
		Quantity   float64 `db:"quantity"`
	}
	var materials []JobMaterial

	err = tx.SelectContext(ctx, &materials, materialQuery, req.JobID)
	if err != nil {
		return fmt.Errorf("failed to get job materials: %w", err)
	}

	// For each material, check if it exists in material_price_log
	for _, material := range materials {
		// Check if material price log exists
		var exists bool
		checkExistsQuery := `
            SELECT EXISTS(
                SELECT 1 
                FROM material_price_log 
                WHERE boq_id = $1 
                AND material_id = $2 
                AND job_id = $3
            )`

		err = tx.GetContext(ctx, &exists, checkExistsQuery, boqID, material.MaterialID, req.JobID)
		if err != nil {
			return fmt.Errorf("failed to check material price log existence: %w", err)
		}

		if !exists {
			insertPriceLogQuery := `
                INSERT INTO material_price_log (
                    material_id, boq_id, job_id, quantity, updated_at
                ) VALUES (
                    $1, $2, $3, $4, CURRENT_TIMESTAMP
                )`

			_, err = tx.ExecContext(ctx, insertPriceLogQuery,
				material.MaterialID,
				boqID,
				req.JobID,
				material.Quantity,
			)
			if err != nil {
				return fmt.Errorf("failed to create material price log: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *boqRepository) UpdateBOQJob(ctx context.Context, boqID uuid.UUID, req requests.BOQJobRequest) error {
	jobID := req.JobID
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check BOQ status
	var status string
	checkStatusQuery := `SELECT status FROM boq WHERE boq_id = $1`
	err = tx.GetContext(ctx, &status, checkStatusQuery, boqID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("boq not found")
		}
		return fmt.Errorf("failed to get BOQ status: %w", err)
	}

	if status != "draft" {
		return errors.New("can only update jobs in BOQ in draft status")
	}

	// Update BOQ job
	updateBOQJobQuery := `
		UPDATE boq_job
		SET quantity = $1, labor_cost = $2
		WHERE boq_id = $3 AND job_id = $4`

	_, err = tx.ExecContext(ctx, updateBOQJobQuery, req.Quantity, req.LaborCost, boqID, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job in BOQ: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *boqRepository) DeleteBOQJob(ctx context.Context, boqID uuid.UUID, jobID uuid.UUID) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check BOQ status
	var status string
	checkStatusQuery := `SELECT status FROM boq WHERE boq_id = $1`
	err = tx.GetContext(ctx, &status, checkStatusQuery, boqID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("boq not found")
		}
		return fmt.Errorf("failed to get BOQ status: %w", err)
	}

	if status != "draft" {
		return errors.New("can only delete jobs from BOQ in draft status")
	}

	// Delete related material price logs first (foreign key constraint)
	deleteMaterialPriceLogQuery := `
        DELETE FROM material_price_log 
        WHERE boq_id = $1 
        AND job_id = $2`

	_, err = tx.ExecContext(ctx, deleteMaterialPriceLogQuery, boqID, jobID)
	if err != nil {
		return fmt.Errorf("failed to delete material price logs: %w", err)
	}

	// Then delete the BOQ job
	deleteBOQJobQuery := `
        DELETE FROM boq_job 
        WHERE boq_id = $1 
        AND job_id = $2`

	result, err := tx.ExecContext(ctx, deleteBOQJobQuery, boqID, jobID)
	if err != nil {
		return fmt.Errorf("failed to delete job from BOQ: %w", err)
	}

	// Check if the BOQ job was actually deleted
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("job not found in BOQ")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
