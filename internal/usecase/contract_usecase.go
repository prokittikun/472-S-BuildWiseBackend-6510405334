package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"time"

	"github.com/google/uuid"
)

type ContractUseCase interface {
	Create(ctx context.Context, req *requests.CreateContractRequest) error
	Update(ctx context.Context, projectID uuid.UUID, req *requests.UpdateContractRequest) error
	Delete(ctx context.Context, projectID uuid.UUID) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*responses.ContractResponse, error)
}

type contractUseCase struct {
	contractRepo repositories.ContractRepository
	projectRepo  repositories.ProjectRepository
}

func NewContractUsecase(
	contractRepo repositories.ContractRepository,
	projectRepo repositories.ProjectRepository,
) ContractUseCase {
	return &contractUseCase{
		contractRepo: contractRepo,
		projectRepo:  projectRepo,
	}
}

func (u *contractUseCase) Create(ctx context.Context, req *requests.CreateContractRequest) error {
	// Validate project exists
	if err := u.projectRepo.ValidateProjectStatus(ctx, req.ProjectID); err != nil {
		return err
	}

	// Create contract model
	contract := &models.Contract{
		ProjectID:           req.ProjectID,
		ProjectDescription:  req.ProjectDescription,
		AreaSize:            req.AreaSize,
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		ForceMajeure:        req.ForceMajeure,
		BreachOfContract:    req.BreachOfContract,
		EndOfContract:       req.EndOfContract,
		TerminationContract: req.TerminationContract,
		Amendment:           req.Amendment,
		GuaranteeWithin:     req.GuaranteeWithin,
		RetentionMoney:      req.RetentionMoney,
		PayWithin:           req.PayWithin,
		ValidateWithin:      req.ValidateWithin,
		Format:              models.StringArray(req.Format),
		CreatedAt:           time.Now(),
	}

	// Create contract
	if err := u.contractRepo.Create(ctx, contract); err != nil {
		return err
	}

	return nil
}

func (u *contractUseCase) Update(ctx context.Context, projectID uuid.UUID, req *requests.UpdateContractRequest) error {
	// Get existing contract
	contract, err := u.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return err
	}

	// Update fields
	if req.ProjectDescription != "" {
		contract.ProjectDescription = req.ProjectDescription
	}
	if req.AreaSize != 0 {
		contract.AreaSize = req.AreaSize
	}
	if !req.StartDate.IsZero() {
		contract.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		contract.EndDate = req.EndDate
	}
	contract.ForceMajeure = req.ForceMajeure
	contract.BreachOfContract = req.BreachOfContract
	contract.EndOfContract = req.EndOfContract
	contract.TerminationContract = req.TerminationContract
	contract.Amendment = req.Amendment
	if req.GuaranteeWithin != 0 {
		contract.GuaranteeWithin = req.GuaranteeWithin
	}
	if req.RetentionMoney != 0 {
		contract.RetentionMoney = req.RetentionMoney
	}
	if req.PayWithin != 0 {
		contract.PayWithin = req.PayWithin
	}
	if req.ValidateWithin != 0 {
		contract.ValidateWithin = req.ValidateWithin
	}
	if len(req.Format) > 0 {
		contract.Format = models.StringArray(req.Format)
	}
	contract.UpdatedAt = time.Now()

	// Update contract
	if err := u.contractRepo.Update(ctx, contract); err != nil {
		return err
	}

	return nil
}

func (u *contractUseCase) Delete(ctx context.Context, projectID uuid.UUID) error {
	return u.contractRepo.Delete(ctx, projectID)
}

func (u *contractUseCase) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*responses.ContractResponse, error) {
	contract, err := u.contractRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	response := &responses.ContractResponse{
		ContractID:          contract.ContractID,
		ProjectID:           contract.ProjectID,
		ProjectDescription:  contract.ProjectDescription,
		AreaSize:            contract.AreaSize,
		StartDate:           contract.StartDate,
		EndDate:             contract.EndDate,
		ForceMajeure:        contract.ForceMajeure,
		BreachOfContract:    contract.BreachOfContract,
		EndOfContract:       contract.EndOfContract,
		TerminationContract: contract.TerminationContract,
		Amendment:           contract.Amendment,
		GuaranteeWithin:     contract.GuaranteeWithin,
		RetentionMoney:      contract.RetentionMoney,
		PayWithin:           contract.PayWithin,
		ValidateWithin:      contract.ValidateWithin,
		Format:              []string(contract.Format),
		CreatedAt:           contract.CreatedAt,
		UpdatedAt:           contract.UpdatedAt,
		Periods:             make([]responses.PeriodResponse, len(contract.Periods)),
	}

	// Map periods and jobs
	for i, period := range contract.Periods {
		periodResponse := responses.PeriodResponse{
			PeriodID:        period.PeriodID,
			PeriodNumber:    period.PeriodNumber,
			AmountPeriod:    period.AmountPeriod,
			DeliveredWithin: period.DeliveredWithin,
			Jobs:            make([]responses.JobPeriodResponse, len(period.Jobs)),
		}

		for j, job := range period.Jobs {
			periodResponse.Jobs[j] = responses.JobPeriodResponse{
				JobID:     job.JobID,
				JobAmount: job.JobAmount,
				Job:       responses.JobResponse{}, // Map job details here
			}
		}

		response.Periods[i] = periodResponse
	}

	return response, nil
}
