package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type ProjectUsecase interface {
	CreateProject(ctx context.Context, req requests.CreateProjectRequest) (*models.Project, error)
	ListProjects(ctx context.Context) ([]*responses.ProjectResponse, error)
	GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req requests.UpdateProjectRequest) (*models.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type projectUsecase struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectUsecase(projectRepo repositories.ProjectRepository) ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

func (pu *projectUsecase) CreateProject(ctx context.Context, req requests.CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		ProjectID:   uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Status:      "planning",
		ContractURL: req.ContractURL,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := pu.projectRepo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pu *projectUsecase) ListProjects(ctx context.Context) ([]*responses.ProjectResponse, error) {
	return pu.projectRepo.ListProjects(ctx)
}

func (pu *projectUsecase) GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return pu.projectRepo.GetProjectByID(ctx, id)
}

func (pu *projectUsecase) UpdateProject(ctx context.Context, id uuid.UUID, req requests.UpdateProjectRequest) (*models.Project, error) {
	project, err := pu.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.ContractURL != "" {
		project.ContractURL = req.ContractURL
	}
	if !req.StartDate.IsZero() {
		project.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		project.EndDate = req.EndDate
	}

	err = pu.projectRepo.UpdateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pu *projectUsecase) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return pu.projectRepo.DeleteProject(ctx, id)
}
