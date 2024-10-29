// usecase/project_usecase.go
package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProjectUsecase interface {
	Create(ctx context.Context, req requests.CreateProjectRequest) (*responses.ProjectResponse, error)
	Update(ctx context.Context, id uuid.UUID, req requests.UpdateProjectRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*responses.ProjectResponse, error)
	List(ctx context.Context) (*responses.ProjectListResponse, error)
	Cancel(ctx context.Context, id uuid.UUID) error

	UpdateProjectStatus(ctx context.Context, req requests.UpdateProjectStatusRequest) error
}

type projectUsecase struct {
	projectRepo repositories.ProjectRepository
	clientRepo  repositories.ClientRepository
}

func NewProjectUsecase(projectRepo repositories.ProjectRepository, clientRepo repositories.ClientRepository) ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
		clientRepo:  clientRepo,
	}
}

func (u *projectUsecase) Create(ctx context.Context, req requests.CreateProjectRequest) (*responses.ProjectResponse, error) {
	client, err := u.clientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	project, err := u.projectRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responses.ProjectResponse{
		ID:          project.ProjectID,
		Name:        project.Name,
		Description: project.Description,
		Address:     project.Address,
		Status:      project.Status,
		ClientID:    project.ClientID,
		Client: &responses.ClientResponse{
			ID:      client.ClientID,
			Name:    client.Name,
			Email:   client.Email,
			Tel:     client.Tel,
			Address: client.Address,
			TaxID:   client.TaxID,
		},
		CreatedAt: project.CreatedAt,
	}, nil
}

func (u *projectUsecase) Update(ctx context.Context, id uuid.UUID, req requests.UpdateProjectRequest) error {
	_, err := u.projectRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("project not found")
	}

	return u.projectRepo.Update(ctx, id, req)

}

func (u *projectUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.projectRepo.Delete(ctx, id)
}

func (u *projectUsecase) GetByID(ctx context.Context, id uuid.UUID) (*responses.ProjectResponse, error) {
	project, client, err := u.projectRepo.GetByIDWithClient(ctx, id)
	if err != nil {
		return nil, err
	}

	return &responses.ProjectResponse{
		ID:          project.ProjectID,
		Name:        project.Name,
		Description: project.Description,
		Address:     project.Address,
		Status:      project.Status,
		ClientID:    project.ClientID,
		Client: &responses.ClientResponse{
			ID:      client.ClientID,
			Name:    client.Name,
			Email:   client.Email,
			Tel:     client.Tel,
			Address: client.Address,
			TaxID:   client.TaxID,
		},
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt.Time,
	}, nil
}

func (u *projectUsecase) List(
	ctx context.Context,
) (*responses.ProjectListResponse, error) {

	projects, err := u.projectRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]responses.ProjectResponse, len(projects))
	for i, project := range projects {
		client, err := u.clientRepo.GetByID(ctx, project.ClientID)
		if err != nil {
			return nil, err
		}

		projectResponses[i] = responses.ProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Address:     project.Address,
			Status:      project.Status,
			ClientID:    project.ClientID,
			Client: &responses.ClientResponse{
				ID:      client.ClientID,
				Name:    client.Name,
				Email:   client.Email,
				Tel:     client.Tel,
				Address: client.Address,
				TaxID:   client.TaxID,
			},
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt.Time,
		}
	}

	return &responses.ProjectListResponse{
		Projects: projectResponses,
	}, nil
}

func (u *projectUsecase) Cancel(ctx context.Context, id uuid.UUID) error {
	return u.projectRepo.Cancel(ctx, id)
}

func (u *projectUsecase) UpdateProjectStatus(ctx context.Context, req requests.UpdateProjectStatusRequest) error {
	return u.projectRepo.UpdateStatus(ctx, req.ProjectID, req.Status)
}
