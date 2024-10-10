package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"time"

	"github.com/google/uuid"
)

type ClientUsecase interface {
	CreateClient(ctx context.Context, req requests.CreateClientRequest) (*models.Client, error)
	ListClients(ctx context.Context) ([]*models.Client, error)
	GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error)
	UpdateClient(ctx context.Context, id uuid.UUID, req requests.UpdateClientRequest) (*models.Client, error)
	DeleteClient(ctx context.Context, id uuid.UUID) error
}

type clientUsecase struct {
	clientRepo repositories.ClientRepository
}

func NewClientUsecase(clientRepo repositories.ClientRepository) ClientUsecase {
	return &clientUsecase{
		clientRepo: clientRepo,
	}
}

func (cu *clientUsecase) CreateClient(ctx context.Context, req requests.CreateClientRequest) (*models.Client, error) {
	client := &models.Client{
		ClientID:      uuid.New(),
		CompanyName:   req.CompanyName,
		ContactPerson: req.ContactPerson,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		TaxID:         req.TaxID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := cu.clientRepo.CreateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (cu *clientUsecase) ListClients(ctx context.Context) ([]*models.Client, error) {
	return cu.clientRepo.ListClients(ctx)
}

func (cu *clientUsecase) GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return cu.clientRepo.GetClientByID(ctx, id)
}

func (cu *clientUsecase) UpdateClient(ctx context.Context, id uuid.UUID, req requests.UpdateClientRequest) (*models.Client, error) {
	client, err := cu.clientRepo.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.CompanyName != "" {
		client.CompanyName = req.CompanyName
	}
	if req.ContactPerson != "" {
		client.ContactPerson = req.ContactPerson
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.TaxID != "" {
		client.TaxID = req.TaxID
	}
	client.UpdatedAt = time.Now()

	err = cu.clientRepo.UpdateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (cu *clientUsecase) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return cu.clientRepo.DeleteClient(ctx, id)
}
