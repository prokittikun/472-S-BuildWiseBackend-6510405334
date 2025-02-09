package rest

import (
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	contractUseCase usecase.ContractUseCase
}

func NewContractHandler(contractUseCase usecase.ContractUseCase) *ContractHandler {
	return &ContractHandler{
		contractUseCase: contractUseCase,
	}
}

func (h *ContractHandler) ContractRoutes(app *fiber.App) {

}
