package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type quotationRepository struct {
	db *sqlx.DB
}

func NewQuotationRepository(db *sqlx.DB) repositories.QuotationRepository {
	return &quotationRepository{db: db}
}

func (r *quotationRepository) CheckBOQStatus(ctx context.Context, projectID uuid.UUID) (string, error) {
	var status string
	query := `
        SELECT b.status 
        FROM boq b 
        WHERE b.project_id = $1`

	err := r.db.GetContext(ctx, &status, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("BOQ not found")
		}
		return "", fmt.Errorf("failed to check BOQ status: %w", err)
	}

	return status, nil
}

func (r *quotationRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.Quotation, error) {
	var quotation models.Quotation
	query := `SELECT * FROM quotation WHERE project_id = $1`

	err := r.db.GetContext(ctx, &quotation, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No quotation exists
		}
		return nil, fmt.Errorf("failed to get quotation: %w", err)
	}

	return &quotation, nil
}

func (r *quotationRepository) Create(ctx context.Context, projectID uuid.UUID) (*models.Quotation, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	quotation := &models.Quotation{
		QuotationID:   uuid.New(),
		ProjectID:     projectID,
		Status:        "draft",
		ValidDate:     sql.NullTime{Time: time.Now().AddDate(0, 1, 0), Valid: true}, // Default validity: 1 month
		TaxPercentage: sql.NullFloat64{Float64: 7, Valid: true},                     // Default tax percentage
	}

	query := `
        INSERT INTO quotation (
            quotation_id, project_id, status, valid_date, 
            final_amount, tax_percentage
        ) VALUES (
            :quotation_id, :project_id, :status, :valid_date, 
            :final_amount, :tax_percentage
        ) RETURNING *`

	_, err = tx.NamedExecContext(ctx, query, quotation)
	if err != nil {
		return nil, fmt.Errorf("failed to create quotation: %w", err)
	}

	err = tx.GetContext(ctx, quotation, "SELECT * FROM quotation WHERE quotation_id = $1", quotation.QuotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created quotation: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return quotation, nil
}
func (r *quotationRepository) GetQuotationJobs(ctx context.Context, projectID uuid.UUID) ([]models.QuotationJob, error) {
	query := `
        SELECT 
            q.quotation_id, 
            q.status, 
            q.valid_date, 
            q.tax_percentage,
            j.name, 
            j.unit, 
            bj.quantity, 
            bj.labor_cost,
            (bj.labor_cost * bj.quantity) as total_labor_cost,
            COALESCE(SUM(mpl.estimated_price), 0) as estimated_price,
            COALESCE(SUM(mpl.estimated_price) * bj.quantity, 0) as total_estimated_price,
            COALESCE((bj.labor_cost * bj.quantity) + (SUM(mpl.estimated_price) * bj.quantity), 
                     bj.labor_cost * bj.quantity) as total,
            bj.selling_price
        FROM project p
        LEFT JOIN quotation q ON q.project_id = p.project_id
        JOIN boq b ON b.project_id = p.project_id
        JOIN boq_job bj ON bj.boq_id = b.boq_id
        JOIN job j ON j.job_id = bj.job_id
        LEFT JOIN material_price_log mpl ON mpl.job_id = j.job_id AND mpl.boq_id = b.boq_id
        WHERE p.project_id = $1
        GROUP BY 
            q.quotation_id, 
            q.status, 
            q.valid_date, 
            q.tax_percentage,
            j.name, 
            j.unit, 
            bj.quantity, 
            bj.labor_cost, 
            bj.selling_price`

	var jobs []models.QuotationJob
	err := r.db.SelectContext(ctx, &jobs, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation jobs: %w", err)
	}

	return jobs, nil
}

func (r *quotationRepository) GetQuotationGeneralCosts(ctx context.Context, projectID uuid.UUID) ([]models.QuotationGeneralCost, error) {
	query := `
        SELECT 
            b.boq_id, gc.g_id, gc.type_name, gc.estimated_cost
        FROM project p
        JOIN boq b ON b.project_id = p.project_id
        LEFT JOIN general_cost gc ON gc.boq_id = b.boq_id
        LEFT JOIN "type" t ON t.type_name = gc.type_name
        WHERE p.project_id = $1`

	var costs []models.QuotationGeneralCost
	err := r.db.SelectContext(ctx, &costs, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation general costs: %w", err)
	}

	return costs, nil
}

func (r *quotationRepository) ValidateApproval(ctx context.Context, projectID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check BOQ status
	var boqStatus string
	boqQuery := `SELECT status FROM boq WHERE project_id = $1`
	err = tx.GetContext(ctx, &boqStatus, boqQuery, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("BOQ not found")
		}
		return fmt.Errorf("failed to get BOQ status: %w", err)
	}

	if boqStatus != "approved" {
		return errors.New("BOQ must be approved before approving quotation")
	}

	// Check quotation status
	var quotationStatus string
	quotationQuery := `SELECT status FROM quotation WHERE project_id = $1`
	err = tx.GetContext(ctx, &quotationStatus, quotationQuery, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("quotation not found")
		}
		return fmt.Errorf("failed to get quotation status: %w", err)
	}

	if quotationStatus != "draft" {
		return errors.New("only draft quotations can be approved")
	}

	return tx.Commit()
}

func (r *quotationRepository) ApproveQuotation(ctx context.Context, projectID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update quotation status
	query := `
        UPDATE quotation 
        SET status = 'approved'
        WHERE project_id = $1 AND status = 'draft'
        RETURNING quotation_id`

	var quotationID uuid.UUID
	err = tx.GetContext(ctx, &quotationID, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no draft quotation found to approve")
		}
		return fmt.Errorf("failed to approve quotation: %w", err)
	}

	return tx.Commit()
}

func (r *quotationRepository) GetQuotationStatus(ctx context.Context, projectID uuid.UUID) (string, error) {
	var status string
	query := `SELECT status FROM quotation WHERE project_id = $1`

	err := r.db.GetContext(ctx, &status, query, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("quotation not found")
		}
		return "", fmt.Errorf("failed to get quotation status: %w", err)
	}

	return status, nil
}

func (r *quotationRepository) GetExportData(ctx context.Context, projectID uuid.UUID) (*models.QuotationExportData, error) {
	// Get main quotation data
	query := `
        SELECT 
            p.project_id,
            p.name,
            p.description,
            p.address,
            c.name as client_name,
            c.address as client_address,
            c.email as client_email,
            c.tel as client_tel,
            c.tax_id as client_tax_id,
            q.quotation_id,
            q.valid_date,
            q.tax_percentage,
            q.final_amount,
            q.status
        FROM project p
        LEFT JOIN client c ON c.client_id = p.client_id
        LEFT JOIN quotation q ON q.project_id = p.project_id
        WHERE p.project_id = $1`

	var data models.QuotationExportData
	err := r.db.GetContext(ctx, &data, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation data: %w", err)
	}

	// Get job details
	jobQuery := `
        SELECT 
            j.name,
            j.unit,
            bj.quantity,
            bj.selling_price,
            (bj.selling_price * bj.quantity) as amount
        FROM boq b
        JOIN boq_job bj ON bj.boq_id = b.boq_id
        JOIN job j ON j.job_id = bj.job_id
        WHERE b.project_id = $1
        ORDER BY j.name`

	err = r.db.SelectContext(ctx, &data.JobDetails, jobQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	return &data, nil
}
