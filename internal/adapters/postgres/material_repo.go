package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type materialRepository struct {
	db *sqlx.DB
}

func NewMaterialRepository(db *sqlx.DB) repositories.MaterialRepository {
	return &materialRepository{
		db: db,
	}
}

func (r *materialRepository) Create(ctx context.Context, req requests.CreateMaterialRequest) (*models.Material, error) {
	material := &models.Material{
		MaterialID: uuid.New().String(),
		Name:       req.Name,
		Unit:       req.Unit,
	}

	query := `
        INSERT INTO Material (
            material_id, name, unit
        ) VALUES (
            :material_id, :name, :unit
        ) RETURNING *`

	rows, err := r.db.NamedQueryContext(ctx, query, material)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("material ID already exists")
		}
		return nil, fmt.Errorf("failed to create material: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(material)
		if err != nil {
			return nil, fmt.Errorf("failed to scan material: %w", err)
		}
		return material, nil
	}
	return nil, errors.New("failed to create material: no rows returned")
}

func (r *materialRepository) Update(ctx context.Context, materialID string, req requests.UpdateMaterialRequest) error {
	query := `
        UPDATE Material SET 
            name = :name,
            unit = :unit
        WHERE material_id = :material_id`

	params := map[string]interface{}{
		"material_id": materialID,
		"name":        req.Name,
		"unit":        req.Unit,
	}

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("material not found")
	}

	return nil
}

func (r *materialRepository) Delete(ctx context.Context, materialID string) error {
	query := `DELETE FROM Material WHERE material_id = $1`

	result, err := r.db.ExecContext(ctx, query, materialID)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return errors.New("material is in use and cannot be deleted")
		}
		return fmt.Errorf("failed to delete material: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("material not found")
	}

	return nil
}

func (r *materialRepository) GetByID(ctx context.Context, materialID string) (*models.Material, error) {
	material := &models.Material{}
	query := `SELECT * FROM Material WHERE material_id = $1`

	err := r.db.GetContext(ctx, material, query, materialID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("material not found")
		}
		return nil, fmt.Errorf("failed to get material: %w", err)
	}

	return material, nil
}

func (r *materialRepository) List(ctx context.Context) ([]models.Material, error) {
	var materials []models.Material
	var args []interface{}

	query := `
		SELECT * FROM Material
	   `

	err := r.db.SelectContext(ctx, &materials, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials: %w", err)
	}

	return materials, nil
}
