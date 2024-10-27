package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GeneralCostHandler struct {
	generalCostUseCase usecase.GeneralCostUseCase
}

func NewGeneralCostHandler(generalCostUseCase usecase.GeneralCostUseCase) *GeneralCostHandler {
	return &GeneralCostHandler{
		generalCostUseCase: generalCostUseCase,
	}
}

func (h *GeneralCostHandler) GeneralCostRoutes(app *fiber.App) {
	generalCost := app.Group("/general-costs")

	generalCost.Post("/", h.Create)
	generalCost.Get("/boq/:boqId", h.GetByBOQID)
	generalCost.Get("/:id", h.GetByID)
	generalCost.Put("/:id", h.Update)
}

func (h *GeneralCostHandler) Create(c *fiber.Ctx) error {
	var req requests.CreateGeneralCostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.BOQID == uuid.Nil || req.TypeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	generalCost, err := h.generalCostUseCase.Create(c.Context(), req)
	if err != nil {
		switch err.Error() {
		case "can only add general cost to BOQ in draft status":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "BOQ must be in draft status",
			})
		case "general cost already exists for this type":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "General cost already exists for this type",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create general cost",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "General cost created successfully",
		"data":    generalCost,
	})
}

func (h *GeneralCostHandler) GetByBOQID(c *fiber.Ctx) error {
	boqID, err := uuid.Parse(c.Params("boqId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid BOQ ID",
		})
	}

	generalCosts, err := h.generalCostUseCase.GetByBOQID(c.Context(), boqID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get general costs",
		})
	}

	return c.JSON(fiber.Map{
		"message": "General costs retrieved successfully",
		"data":    generalCosts,
	})
}

func (h *GeneralCostHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid general cost ID",
		})
	}

	generalCost, err := h.generalCostUseCase.GetByID(c.Context(), id)
	if err != nil {
		switch err.Error() {
		case "general cost not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "General cost not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get general cost",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "General cost retrieved successfully",
		"data":    generalCost,
	})
}

func (h *GeneralCostHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid general cost ID",
		})
	}

	var req requests.UpdateGeneralCostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate estimated cost
	if req.EstimatedCost < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Estimated cost must be positive",
		})
	}

	err = h.generalCostUseCase.Update(c.Context(), id, req)
	if err != nil {
		switch err.Error() {
		case "general cost not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "General cost not found",
			})
		case "can only update general cost for BOQ in draft status":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "BOQ must be in draft status",
			})
		case "estimated cost must be positive":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Estimated cost must be positive",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update general cost",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "General cost updated successfully",
	})
}
