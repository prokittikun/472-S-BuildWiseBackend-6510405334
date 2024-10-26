package rest

import (
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BOQHandler struct {
	boqUsecase usecase.BOQUsecase
}

func NewBOQHandler(boqUsecase usecase.BOQUsecase) *BOQHandler {
	return &BOQHandler{
		boqUsecase: boqUsecase,
	}
}

func (h *BOQHandler) BOQRoutes(app *fiber.App) {
	boq := app.Group("/boqs")

	boq.Get("/project/:project_id", h.GetBoqWithProject)
}

func (h *BOQHandler) GetBoqWithProject(c *fiber.Ctx) error {
	project_id := c.Params("project_id")
	if project_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	uuid, err := uuid.Parse(project_id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	boq, err := h.boqUsecase.GetBoqWithProject(c.Context(), uuid)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "BOQ retrieved successfully",
		"data":    boq,
	})
}
