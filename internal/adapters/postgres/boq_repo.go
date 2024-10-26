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

func (r *boqRepository) Create(ctx context.Context, req requests.CreateBOQRequest) (*models.BOQ, error) {
	boq := &models.BOQ{
		BOQID:     uuid.New(),
		ProjectID: req.ProjectID,
		Status:    "draft",
	}

	query := `
        INSERT INTO BOQ (
            boq_id, project_id, status
        ) VALUES (
            :boq_id, :project_id, :status
        ) RETURNING *`

	rows, err := r.db.NamedQueryContext(ctx, query, boq)
	if err != nil {
		return nil, fmt.Errorf("failed to create BOQ: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(boq)
		if err != nil {
			return nil, fmt.Errorf("failed to scan BOQ: %w", err)
		}
		return boq, nil
	}
	return nil, errors.New("failed to create BOQ: no rows returned")
}

func (r *boqRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BOQ, error) {
	boq := &models.BOQ{}
	query := `SELECT * FROM BOQ WHERE boq_id = $1`

	err := r.db.GetContext(ctx, boq, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("BOQ not found")
		}
		return nil, fmt.Errorf("failed to get BOQ: %w", err)
	}

	return boq, nil
}

func (r *boqRepository) GetByIDWithProject(ctx context.Context, id uuid.UUID) (*models.BOQ, *models.Project, error) {
	boq := &models.BOQ{}
	project := &models.Project{}

	query := `
        SELECT 
            b.*,
            p.project_id as "project.project_id",
            p.name as "project.name",
            p.description as "project.description",
            p.address as "project.address",
            p.status as "project.status",
            p.client_id as "project.client_id",
            p.created_at as "project.created_at",
            p.updated_at as "project.updated_at"
        FROM BOQ b
        LEFT JOIN Project p ON b.project_id = p.project_id
        WHERE b.boq_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get BOQ with project: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&boq.BOQID, &boq.ProjectID, &boq.Status, &boq.SellingGeneralCost,
			&project.ProjectID, &project.Name, &project.Description,
			&project.Address, &project.Status, &project.ClientID,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan BOQ with project: %w", err)
		}
		return boq, project, nil
	}

	return nil, nil, errors.New("BOQ not found")
}

func (r *boqRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.BOQ, error) {
	boq := &models.BOQ{}
	query := `SELECT * FROM BOQ WHERE project_id = $1`

	err := r.db.GetContext(ctx, boq, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("BOQ not found")
		}
		return nil, fmt.Errorf("failed to get BOQ: %w", err)
	}

	return boq, nil
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
