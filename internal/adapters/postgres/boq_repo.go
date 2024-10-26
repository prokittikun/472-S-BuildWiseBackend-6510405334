package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
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
