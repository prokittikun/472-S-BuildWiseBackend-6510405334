package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/responses"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type projectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) repositories.ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

func (pr *projectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO Project (project_id, name, description, status, contract_url, start_date, end_date, created_at, updated_at)
			  VALUES (:project_id, :name, :description, :status, :contract_url, :start_date, :end_date, :created_at, :updated_at)`

	_, err := pr.db.NamedExecContext(ctx, query, project)
	return err
}

func (pr *projectRepository) ListProjects(ctx context.Context) ([]*responses.ProjectResponse, error) {
	query := `
        SELECT 
          *
        FROM 
            Project p
        JOIN 
            Client c ON p.client_id = c.client_id
    `

	rows, err := pr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*responses.ProjectResponse

	for rows.Next() {
		var p responses.ProjectResponse
		var c responses.ClientResponse
		var quotationID, contractID, invoiceID, bID, clientID sql.NullString

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Status, &p.ContractURL,
			&p.StartDate, &p.EndDate, &quotationID, &contractID,
			&invoiceID, &bID, &clientID, &p.CreatedAt, &p.UpdatedAt,
			&c.ClientID, &c.CompanyName, &c.ContactPerson, &c.Email,
			&c.Phone, &c.Address, &c.TaxID, &c.CreatedAt, &c.UpdatedAt,
		)
		fmt.Print(err)
		if err != nil {
			return nil, err
		}

		p.QuotationID = nullStringToNullableUUID(quotationID)
		p.ContractID = nullStringToNullableUUID(contractID)
		p.InvoiceID = nullStringToNullableUUID(invoiceID)
		p.BID = nullStringToNullableUUID(bID)
		p.ClientID = nullStringToNullableUUID(clientID)

		p.Client = c
		projects = append(projects, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func nullStringToNullableUUID(ns sql.NullString) responses.NullableUUID {
	if !ns.Valid {
		return responses.NullableUUID{Valid: false}
	}
	id, err := uuid.Parse(ns.String)
	if err != nil {
		return responses.NullableUUID{Valid: false}
	}
	return responses.NullableUUID{UUID: id, Valid: true}
}

func (pr *projectRepository) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := pr.db.GetContext(ctx, &project, `SELECT * FROM Project WHERE project_id = $1`, id)
	return &project, err
}

func (pr *projectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `UPDATE Project SET name = :name, description = :description, status = :status, contract_url = :contract_url,
			  start_date = :start_date, end_date = :end_date, updated_at = :updated_at
			  WHERE project_id = :project_id`

	_, err := pr.db.NamedExecContext(ctx, query, project)
	return err
}

func (pr *projectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	_, err := pr.db.ExecContext(ctx, `DELETE FROM Project WHERE project_id = $1`, id)
	return err
}
