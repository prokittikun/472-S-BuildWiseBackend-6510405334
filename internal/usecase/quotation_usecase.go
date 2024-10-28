package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/responses"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type QuotationUsecase interface {
	CreateOrGetQuotation(ctx context.Context, projectID uuid.UUID) (*responses.QuotationResponse, error)
	ApproveQuotation(ctx context.Context, projectID uuid.UUID) error
	ExportQuotation(ctx context.Context, projectID uuid.UUID) (*models.QuotationExportData, error)
}

type quotationUsecase struct {
	quotationRepo repositories.QuotationRepository
}

func NewQuotationUsecase(quotationRepo repositories.QuotationRepository) QuotationUsecase {
	return &quotationUsecase{
		quotationRepo: quotationRepo,
	}
}
func (u *quotationUsecase) buildQuotationResponse(
	quotation *models.Quotation,
	jobs []models.QuotationJob,
	costs []models.QuotationGeneralCost,
) *responses.QuotationResponse {
	response := &responses.QuotationResponse{
		QuotationID: quotation.QuotationID,
		Status:      string(quotation.Status),
		ValidDate:   getValidTime(quotation.ValidDate),
		Jobs:        make([]responses.QuotationJobDetail, 0),
		Costs:       make([]responses.GeneralCostDetail, 0),
	}

	// Process jobs
	var totalLaborCost float64
	var totalMaterialCost float64

	for _, job := range jobs {

		jobDetail := responses.QuotationJobDetail{
			Name:               job.JobName,
			Unit:               job.Unit,
			Quantity:           job.Quantity,
			LaborCost:          job.LaborCost,
			SellingPrice:       job.SellingPrice.Float64,
			TotalMaterialPrice: job.TotalMaterialPrice.Float64,
			Total:              job.Total.Float64,
			OverallCost:        job.OverallCost.Float64,
			TotalSellingPrice:  job.TotalSellingPrice.Float64,
		}

		if job.SellingPrice.Valid {
			jobDetail.SellingPrice = job.SellingPrice.Float64
		}

		// Accumulate totals
		totalLaborCost += job.LaborCost * job.Quantity
		if job.TotalMaterialPrice.Valid {
			totalMaterialCost += job.TotalMaterialPrice.Float64 * job.Quantity
		}

		response.Jobs = append(response.Jobs, jobDetail)
	}

	// Process general costs
	var totalGeneralCost float64
	for _, cost := range costs {
		if cost.EstimatedCost.Valid {
			costDetail := responses.GeneralCostDetail{
				TypeName:      cost.TypeName,
				EstimatedCost: cost.EstimatedCost.Float64,
			}
			totalGeneralCost += cost.EstimatedCost.Float64
			response.Costs = append(response.Costs, costDetail)
		}
	}

	return response
}
func getValidTime(nullTime sql.NullTime) time.Time {
	if nullTime.Valid {
		return nullTime.Time
	}
	return time.Time{}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (u *quotationUsecase) CreateOrGetQuotation(ctx context.Context, projectID uuid.UUID) (*responses.QuotationResponse, error) {
	// Check BOQ status
	boqStatus, err := u.quotationRepo.CheckBOQStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if boqStatus != "approved" {
		return nil, errors.New("BOQ must be approved before creating quotation")
	}

	// Check existing quotation
	quotation, err := u.quotationRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create new quotation if none exists
	if quotation == nil {
		quotation, err = u.quotationRepo.Create(ctx, projectID)
		if err != nil {
			return nil, err
		}
	}

	// Get jobs and costs
	jobs, err := u.quotationRepo.GetQuotationJobs(ctx, projectID)
	if err != nil {
		return nil, err
	}

	costs, err := u.quotationRepo.GetQuotationGeneralCosts(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Build response
	response := u.buildQuotationResponse(quotation, jobs, costs)
	return response, nil
}

func (u *quotationUsecase) ApproveQuotation(ctx context.Context, projectID uuid.UUID) error {
	// Validate approval conditions
	err := u.quotationRepo.ValidateApproval(ctx, projectID)
	if err != nil {
		return err
	}

	// If validation passes, approve the quotation
	err = u.quotationRepo.ApproveQuotation(ctx, projectID)
	if err != nil {
		return err
	}

	// Get updated quotation details for response
	quotation, err := u.quotationRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get updated quotation: %w", err)
	}

	jobs, err := u.quotationRepo.GetQuotationJobs(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get quotation jobs: %w", err)
	}

	costs, err := u.quotationRepo.GetQuotationGeneralCosts(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get quotation costs: %w", err)
	}

	// Build and return response
	_ = u.buildQuotationResponse(quotation, jobs, costs)
	return nil
}

func (u *quotationUsecase) ExportQuotation(ctx context.Context, projectID uuid.UUID) (*models.QuotationExportData, error) {

	// Check BOQ status
	boqStatus, err := u.quotationRepo.CheckBOQStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if boqStatus != "approved" {
		return nil, errors.New("BOQ must be approved before exporting quotation")
	}

	quotationStatus, err := u.quotationRepo.GetQuotationStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if quotationStatus != "approved" {
		return nil, errors.New("only approved quotations can be exported")
	}

	// Get export data
	exportData, err := u.quotationRepo.GetExportData(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return exportData, nil

}
