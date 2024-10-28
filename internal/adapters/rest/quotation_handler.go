// rest/quotation_handler.go
package rest

import (
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type QuotationHandler struct {
	quotationUsecase usecase.QuotationUsecase
}

func NewQuotationHandler(quotationUsecase usecase.QuotationUsecase) *QuotationHandler {
	return &QuotationHandler{
		quotationUsecase: quotationUsecase,
	}
}

func (h *QuotationHandler) QuotationRoutes(app *fiber.App) {
	quotation := app.Group("/quotations")

	// Create or Get Quotation

	quotation.Post("/projects/:projectId", h.CreateOrGetQuotation)
	quotation.Put("/projects/:projectId/approve", h.ApproveQuotation)

}

// CreateOrGetQuotation creates a new quotation or retrieves existing one
func (h *QuotationHandler) CreateOrGetQuotation(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	response, err := h.quotationUsecase.CreateOrGetQuotation(c.Context(), projectID)
	if err != nil {
		switch err.Error() {
		case "BOQ must be approved before creating quotation":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "BOQ not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "BOQ not found for this project",
			})
		case "project not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Project not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to process quotation",
				"details": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Quotation processed successfully",
		"data":    response,
	})
}

func (h *QuotationHandler) ApproveQuotation(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	err = h.quotationUsecase.ApproveQuotation(c.Context(), projectID)
	if err != nil {
		switch err.Error() {
		case "BOQ not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "BOQ not found",
			})
		case "quotation not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Quotation not found",
			})
		case "BOQ must be approved before approving quotation":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "only draft quotations can be approved":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "no draft quotation found to approve":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No draft quotation found to approve",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to approve quotation",
				"details": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Quotation approved successfully",
	})
}
