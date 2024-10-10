package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientHandler struct {
	clientUsecase usecase.ClientUsecase
}

func NewClientHandler(clientUsecase usecase.ClientUsecase) *ClientHandler {
	return &ClientHandler{
		clientUsecase: clientUsecase,
	}
}

func (h *ClientHandler) ClientRoutes(app *fiber.App) {
	app.Post("/clients", h.CreateClient)
	app.Get("/clients", h.ListClients)
	app.Get("/clients/:id", h.GetClient)
	app.Put("/clients/:id", h.UpdateClient)
	app.Delete("/clients/:id", h.DeleteClient)
}

func (h *ClientHandler) CreateClient(c *fiber.Ctx) error {
	var req requests.CreateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	client, err := h.clientUsecase.CreateClient(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create client"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Client created successfully",
		"data":    responses.CreateClientResponse(*client),
	})
}

func (h *ClientHandler) ListClients(c *fiber.Ctx) error {
	clients, err := h.clientUsecase.ListClients(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch clients"})
	}

	clientsResponse := make([]responses.ClientResponse, len(clients))
	for i, client := range clients {
		clientsResponse[i] = responses.ClientResponse(*client)
	}

	return c.JSON(fiber.Map{
		"data": clientsResponse,
	})
}

func (h *ClientHandler) GetClient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid client ID"})
	}

	client, err := h.clientUsecase.GetClient(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client not found"})
	}

	return c.JSON(fiber.Map{
		"data": responses.ClientResponse(*client),
	})
}

func (h *ClientHandler) UpdateClient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid client ID"})
	}

	var req requests.UpdateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	client, err := h.clientUsecase.UpdateClient(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update client"})
	}

	return c.JSON(fiber.Map{
		"message": "Client updated successfully",
		"data":    responses.UpdateClientResponse(*client),
	})
}

func (h *ClientHandler) DeleteClient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid client ID"})
	}

	err = h.clientUsecase.DeleteClient(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete client"})
	}

	return c.JSON(fiber.Map{
		"message": "Client deleted successfully",
	})
}
