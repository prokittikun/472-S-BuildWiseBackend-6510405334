package usecase

import (
	"boonkosang/internal/repositories"
)

type ContractUseCase interface {
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
