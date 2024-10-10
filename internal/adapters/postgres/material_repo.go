package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"

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

func (mr *materialRepository) CreateMaterial(ctx context.Context, material *models.Material) error {
	query := `INSERT INTO Material (name, type, unit_of_measure, created_at, updated_at)
			  VALUES (:name, :type, :unit_of_measure, :created_at, :updated_at)`

	_, err := mr.db.NamedExecContext(ctx, query, material)
	return err
}

func (mr *materialRepository) ListMaterials(ctx context.Context) ([]*models.Material, error) {
	var materials []*models.Material
	err := mr.db.SelectContext(ctx, &materials, `SELECT * FROM Material`)
	return materials, err
}

func (mr *materialRepository) GetMaterialByName(ctx context.Context, name string) (*models.Material, error) {
	var material models.Material
	err := mr.db.GetContext(ctx, &material, `SELECT * FROM Material WHERE name = $1`, name)
	return &material, err
}

func (mr *materialRepository) UpdateMaterial(ctx context.Context, material *models.Material) error {
	query := `UPDATE Material SET type = :type, unit_of_measure = :unit_of_measure, updated_at = :updated_at
			  WHERE name = :name`

	_, err := mr.db.NamedExecContext(ctx, query, material)
	return err
}

func (mr *materialRepository) DeleteMaterial(ctx context.Context, name string) error {
	_, err := mr.db.ExecContext(ctx, `DELETE FROM Material WHERE name = $1`, name)
	return err
}

func (mr *materialRepository) GetMaterialPriceHistory(ctx context.Context, name string) ([]*models.MaterialPriceLog, error) {
	var priceHistory []*models.MaterialPriceLog
	query := `SELECT * FROM Material_price_log WHERE material_name = $1 ORDER BY created_at DESC`
	err := mr.db.SelectContext(ctx, &priceHistory, query, name)
	return priceHistory, err
}

func (mr *materialRepository) MaterialExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Material WHERE name = $1)"
	err := mr.db.GetContext(ctx, &exists, query, name)
	return exists, err
}
