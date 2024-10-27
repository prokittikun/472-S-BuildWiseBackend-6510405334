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

type generalCostRepository struct {
	db *sqlx.DB
}

func NewGeneralCostRepository(db *sqlx.DB) repositories.GeneralCostRepository {
	return &generalCostRepository{db: db}
}

// Create new general cost
func (r *generalCostRepository) Create(ctx context.Context, generalCost *models.GeneralCost) (*responses.GeneralCostResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if exists
	var exists bool
	checkQuery := `
       SELECT EXISTS(
           SELECT 1 FROM general_cost 
           WHERE boq_id = $1 AND type_name = $2
       )`

	err = tx.GetContext(ctx, &exists, checkQuery, generalCost.BOQID, generalCost.TypeName)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing general cost: %w", err)
	}

	if exists {
		return nil, errors.New("general cost already exists for this type")
	}

	// Insert new general cost
	query := `
       INSERT INTO general_cost (
           g_id, boq_id, type_name, actual_cost, estimated_cost
       ) VALUES (
           $1, $2, $3, $4, $5
       ) RETURNING *`

	var result models.GeneralCost
	err = tx.GetContext(ctx, &result, query,
		generalCost.GID,
		generalCost.BOQID,
		generalCost.TypeName,
		generalCost.ActualCost,
		generalCost.EstimatedCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create general cost: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &responses.GeneralCostResponse{
		GID:           result.GID,
		BOQID:         result.BOQID,
		TypeName:      result.TypeName,
		ActualCost:    result.ActualCost.Float64,
		EstimatedCost: result.EstimatedCost.Float64,
	}, nil
}

// Get general cost by BOQ ID
func (r *generalCostRepository) GetByBOQID(ctx context.Context, boqID uuid.UUID) (*responses.GeneralCostListResponse, error) {
	query := `
       SELECT 
           g_id,
           boq_id,
           type_name,
           actual_cost,
           estimated_cost
       FROM general_cost
       WHERE boq_id = $1`

	var generalCosts []models.GeneralCost
	err := r.db.SelectContext(ctx, &generalCosts, query, boqID)
	if err != nil {
		return nil, fmt.Errorf("failed to get general costs: %w", err)
	}

	var response []responses.GeneralCostResponse
	for _, gc := range generalCosts {
		response = append(response, responses.GeneralCostResponse{
			GID:           gc.GID,
			BOQID:         gc.BOQID,
			TypeName:      gc.TypeName,
			ActualCost:    gc.ActualCost.Float64,
			EstimatedCost: gc.EstimatedCost.Float64,
		})
	}

	return &responses.GeneralCostListResponse{
		GeneralCosts: response,
	}, nil
}

// Get general cost by ID
func (r *generalCostRepository) GetByID(ctx context.Context, gID uuid.UUID) (*models.GeneralCost, error) {
	query := `
       SELECT 
           g_id,
           boq_id,
           type_name,
           actual_cost,
           estimated_cost
       FROM general_cost
       WHERE g_id = $1`

	var generalCost models.GeneralCost
	err := r.db.GetContext(ctx, &generalCost, query, gID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("general cost not found")
		}
		return nil, fmt.Errorf("failed to get general cost: %w", err)
	}

	return &generalCost, nil
}

// Update general cost
func (r *generalCostRepository) Update(ctx context.Context, gID uuid.UUID, req requests.UpdateGeneralCostRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check BOQ status first
	var boqStatus string
	statusQuery := `
       SELECT b.status
       FROM general_cost gc
       JOIN boq b ON b.boq_id = gc.boq_id
       WHERE gc.g_id = $1`

	err = tx.GetContext(ctx, &boqStatus, statusQuery, gID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("general cost not found")
		}
		return fmt.Errorf("failed to check BOQ status: %w", err)
	}

	if boqStatus != "draft" {
		return errors.New("can only update general cost for BOQ in draft status")
	}

	// Validate estimated cost
	if req.EstimatedCost < 0 {
		return errors.New("estimated cost must be positive")
	}

	// Update general cost
	updateQuery := `
       UPDATE general_cost 
       SET estimated_cost = $1
       WHERE g_id = $2`

	result, err := tx.ExecContext(ctx, updateQuery, req.EstimatedCost, gID)
	if err != nil {
		return fmt.Errorf("failed to update general cost: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("general cost not found")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
