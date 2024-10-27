package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type GeneralCostUseCase interface {
	Create(ctx context.Context, req requests.CreateGeneralCostRequest) (*responses.GeneralCostResponse, error)
	GetByBOQID(ctx context.Context, boqID uuid.UUID) (*responses.GeneralCostListResponse, error)
	GetByID(ctx context.Context, gID uuid.UUID) (*responses.GeneralCostResponse, error)
	Update(ctx context.Context, gID uuid.UUID, req requests.UpdateGeneralCostRequest) error
}

type generalCostUseCase struct {
	generalCostRepo repositories.GeneralCostRepository
	boqRepo         repositories.BOQRepository
}

func NewGeneralCostUsecase(generalCostRepo repositories.GeneralCostRepository, boqRepo repositories.BOQRepository) GeneralCostUseCase {
	return &generalCostUseCase{
		generalCostRepo: generalCostRepo,
		boqRepo:         boqRepo,
	}
}

func (u *generalCostUseCase) Create(ctx context.Context, req requests.CreateGeneralCostRequest) (*responses.GeneralCostResponse, error) {
	// Check BOQ status
	boq, err := u.boqRepo.GetByID(ctx, req.BOQID)
	if err != nil {
		return nil, err
	}
	if boq.Status != "draft" {
		return nil, errors.New("can only add general cost to BOQ in draft status")
	}

	// Create general cost model
	generalCost := &models.GeneralCost{
		GID:           uuid.New(),
		BOQID:         req.BOQID,
		TypeName:      req.TypeName,
		ActualCost:    sql.NullFloat64{Float64: 0, Valid: true},
		EstimatedCost: sql.NullFloat64{Float64: 0, Valid: true},
	}

	return u.generalCostRepo.Create(ctx, generalCost)
}

func (u *generalCostUseCase) GetByBOQID(ctx context.Context, boqID uuid.UUID) (*responses.GeneralCostListResponse, error) {
	return u.generalCostRepo.GetByBOQID(ctx, boqID)
}

func (u *generalCostUseCase) GetByID(ctx context.Context, gID uuid.UUID) (*responses.GeneralCostResponse, error) {
	generalCost, err := u.generalCostRepo.GetByID(ctx, gID)
	if err != nil {
		return nil, err
	}

	return &responses.GeneralCostResponse{
		GID:           generalCost.GID,
		BOQID:         generalCost.BOQID,
		TypeName:      generalCost.TypeName,
		ActualCost:    generalCost.ActualCost.Float64,
		EstimatedCost: generalCost.EstimatedCost.Float64,
	}, nil
}

func (u *generalCostUseCase) Update(ctx context.Context, gID uuid.UUID, req requests.UpdateGeneralCostRequest) error {
	// Get general cost to check if exists
	generalCost, err := u.generalCostRepo.GetByID(ctx, gID)
	if err != nil {
		return err
	}

	// Check BOQ status
	boq, err := u.boqRepo.GetByID(ctx, generalCost.BOQID)
	if err != nil {
		return err
	}
	if boq.Status != "draft" {
		return errors.New("can only update general cost for BOQ in draft status")
	}

	// Validate estimated cost
	if req.EstimatedCost < 0 {
		return errors.New("estimated cost must be positive")
	}

	return u.generalCostRepo.Update(ctx, gID, req)
}
