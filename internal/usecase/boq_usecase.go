package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"

	"github.com/google/uuid"
)

type BOQUsecase interface {
	Create(ctx context.Context, req requests.CreateBOQRequest) (*responses.BOQResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*responses.BOQResponse, error)
}

type boqUsecase struct {
	boqRepo     repositories.BOQRepository
	projectRepo repositories.ProjectRepository
}

func NewBOQUsecase(boqRepo repositories.BOQRepository, projectRepo repositories.ProjectRepository) BOQUsecase {
	return &boqUsecase{
		boqRepo:     boqRepo,
		projectRepo: projectRepo,
	}
}

func (u *boqUsecase) Create(ctx context.Context, req requests.CreateBOQRequest) (*responses.BOQResponse, error) {
	// Verify project exists
	project, err := u.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Check if BOQ already exists for this project
	existing, err := u.boqRepo.GetByProjectID(ctx, req.ProjectID)
	if err == nil && existing != nil {
		return nil, errors.New("BOQ already exists for this project")
	}

	boq, err := u.boqRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responses.BOQResponse{
		ID:                 boq.BOQID,
		ProjectID:          boq.ProjectID,
		Status:             boq.Status,
		SellingGeneralCost: boq.SellingGeneralCost.Float64,
		Project: responses.ProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			ClientID:    project.ClientID,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt.Time,
		},
	}, nil
}
func (u *boqUsecase) GetByID(ctx context.Context, id uuid.UUID) (*responses.BOQResponse, error) {
	boq, err := u.boqRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := u.projectRepo.GetByID(ctx, boq.ProjectID)
	if err != nil {
		return nil, err
	}

	return &responses.BOQResponse{
		ID:                 boq.BOQID,
		ProjectID:          boq.ProjectID,
		Status:             boq.Status,
		SellingGeneralCost: boq.SellingGeneralCost.Float64,
		Project: responses.ProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			ClientID:    project.ClientID,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt.Time,
		},
	}, nil
}
