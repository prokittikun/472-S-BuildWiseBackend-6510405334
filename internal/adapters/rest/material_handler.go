package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type MaterialHandler struct {
	materialUsecase usecase.MaterialUsecase
}

func NewMaterialHandler(materialUsecase usecase.MaterialUsecase) *MaterialHandler {
	return &MaterialHandler{
		materialUsecase: materialUsecase,
	}
}

func (h *MaterialHandler) MaterialRoutes(app *fiber.App) {
	app.Post("/materials", h.CreateMaterial)
	app.Get("/materials", h.ListMaterials)
	app.Get("/materials/:name", h.GetMaterial)
	app.Put("/materials/:name", h.UpdateMaterial)
	app.Delete("/materials/:name", h.DeleteMaterial)
	app.Get("/materials/:name/price-history", h.GetMaterialPriceHistory)
}

func (h *MaterialHandler) CreateMaterial(c *fiber.Ctx) error {
	var req requests.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	material, err := h.materialUsecase.CreateMaterial(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create material"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Material created successfully",
		"data":    responses.CreateMaterialResponse(*material),
	})
}

func (h *MaterialHandler) ListMaterials(c *fiber.Ctx) error {
	materials, err := h.materialUsecase.ListMaterials(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch materials"})
	}

	materialsResponse := make([]responses.MaterialResponse, len(materials))
	for i, material := range materials {
		materialsResponse[i] = responses.MaterialResponse(*material)
	}

	return c.JSON(fiber.Map{
		"data": materialsResponse,
	})
}

func (h *MaterialHandler) GetMaterial(c *fiber.Ctx) error {
	name := c.Params("name")

	material, err := h.materialUsecase.GetMaterial(c.Context(), name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Material not found"})
	}

	return c.JSON(fiber.Map{
		"data": responses.MaterialResponse(*material),
	})
}

func (h *MaterialHandler) UpdateMaterial(c *fiber.Ctx) error {
	name := c.Params("name")

	var req requests.UpdateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	material, err := h.materialUsecase.UpdateMaterial(c.Context(), name, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update material"})
	}

	return c.JSON(fiber.Map{
		"message": "Material updated successfully",
		"data":    responses.UpdateMaterialResponse(*material),
	})
}

func (h *MaterialHandler) DeleteMaterial(c *fiber.Ctx) error {
	name := c.Params("name")

	err := h.materialUsecase.DeleteMaterial(c.Context(), name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete material"})
	}

	return c.JSON(fiber.Map{
		"message": "Material deleted successfully",
	})
}

func (h *MaterialHandler) GetMaterialPriceHistory(c *fiber.Ctx) error {
	name := c.Params("name")

	priceHistory, err := h.materialUsecase.GetMaterialPriceHistory(c.Context(), name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch price history"})
	}

	return c.JSON(fiber.Map{
		"data": priceHistory,
	})
}
