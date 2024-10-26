package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type BOQUsecase interface {
	GetBoqWithProject(ctx context.Context, project_id uuid.UUID) (*responses.BOQResponse, error)
	AddBOQJob(ctx context.Context, boqID uuid.UUID, req requests.BOQJobRequest) error
	DeleteBOQJob(ctx context.Context, boqID uuid.UUID, jobID uuid.UUID) error
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
func (u *boqUsecase) GetBoqWithProject(ctx context.Context, project_id uuid.UUID) (*responses.BOQResponse, error) {
	return u.boqRepo.GetBoqWithProject(ctx, project_id)
}

func (u *boqUsecase) AddBOQJob(ctx context.Context, boqID uuid.UUID, req requests.BOQJobRequest) error {
	return u.boqRepo.AddBOQJob(ctx, boqID, req)
}

func (u *boqUsecase) DeleteBOQJob(ctx context.Context, boqID uuid.UUID, jobID uuid.UUID) error {
	return u.boqRepo.DeleteBOQJob(ctx, boqID, jobID)
}
