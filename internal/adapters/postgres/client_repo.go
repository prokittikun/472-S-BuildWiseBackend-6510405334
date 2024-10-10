package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type clientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) repositories.ClientRepository {
	return &clientRepository{
		db: db,
	}
}

func (cr *clientRepository) CreateClient(ctx context.Context, client *models.Client) error {
	query := `INSERT INTO Client (client_id, company_name, contact_person, email, phone, address, tax_id, created_at, updated_at)
			  VALUES (:client_id, :company_name, :contact_person, :email, :phone, :address, :tax_id, :created_at, :updated_at)`

	_, err := cr.db.NamedExecContext(ctx, query, client)
	return err
}

func (cr *clientRepository) ListClients(ctx context.Context) ([]*models.Client, error) {
	var clients []*models.Client
	err := cr.db.SelectContext(ctx, &clients, `SELECT * FROM Client`)
	return clients, err
}

func (cr *clientRepository) GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	var client models.Client
	err := cr.db.GetContext(ctx, &client, `SELECT * FROM Client WHERE client_id = $1`, id)
	return &client, err
}

func (cr *clientRepository) UpdateClient(ctx context.Context, client *models.Client) error {
	query := `UPDATE Client SET company_name = :company_name, contact_person = :contact_person, email = :email,
			  phone = :phone, address = :address, tax_id = :tax_id, updated_at = :updated_at
			  WHERE client_id = :client_id`

	_, err := cr.db.NamedExecContext(ctx, query, client)
	return err
}

func (cr *clientRepository) DeleteClient(ctx context.Context, id uuid.UUID) error {
	_, err := cr.db.ExecContext(ctx, `DELETE FROM Client WHERE client_id = $1`, id)
	return err
}
