package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type contractRepository struct {
	db *sqlx.DB
}

func NewContractRepository(db *sqlx.DB) repositories.ContractRepository {
	return &contractRepository{db: db}
}

// Create ...
func (r *contractRepository) Create(ctx context.Context, contract *models.Contract) error {
	return nil
}

// Update ...
func (r *contractRepository) Update(ctx context.Context, contract *models.Contract) error {
	return nil
}

// Delete ...
func (r *contractRepository) Delete(ctx context.Context, projectID uuid.UUID) error {
	return nil
}

// GetByProjectID ...
func (r *contractRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.Contract, error) {
	return nil, nil
}

// ValidateProjectStatus ...
func (r *contractRepository) ValidateProjectStatus(ctx context.Context, projectID uuid.UUID) error {
	return nil
}
