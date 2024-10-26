package rest

import (
	"boonkosang/internal/requests"
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

	boq.Post("/", h.Create)
	boq.Get("/:id", h.GetByID)
	boq.Get("/project/:project_id", h.GetBoqWithProject)
}

func (h *BOQHandler) Create(c *fiber.Ctx) error {
	var req requests.CreateBOQRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	boq, err := h.boqUsecase.Create(c.Context(), req)
	if err != nil {
		if err.Error() == "project not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "BOQ already exists for this project" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "BOQ created successfully",
		"data":    boq,
	})
}

func (h *BOQHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid BOQ ID",
		})
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid BOQ ID format",
		})
	}

	boq, err := h.boqUsecase.GetByID(c.Context(), uuid)
	if err != nil {
		if err.Error() == "BOQ not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve BOQ",
		})
	}

	return c.JSON(fiber.Map{
		"message": "BOQ retrieved successfully",
		"data":    boq,
	})
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
