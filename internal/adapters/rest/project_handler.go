package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectUsecase usecase.ProjectUsecase
}

func NewProjectHandler(projectUsecase usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: projectUsecase,
	}
}

func (h *ProjectHandler) ProjectRoutes(app *fiber.App) {
	app.Post("/projects", h.CreateProject)
	app.Get("/projects", h.ListProjects)
	app.Get("/projects/:id", h.GetProject)
	app.Put("/projects/:id", h.UpdateProject)
	app.Delete("/projects/:id", h.DeleteProject)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var req requests.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	project, err := h.projectUsecase.CreateProject(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Project created successfully",
		"data": responses.CreateProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			ContractURL: project.ContractURL,
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
		},
	})
}

func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	projects, err := h.projectUsecase.ListProjects(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch projects"})
	}
	// projectsResponse := make([]responses.ProjectResponse, 0)
	// for _, project := range projects {
	// 	projectsResponse = append(projectsResponse, responses.ProjectResponse{
	// 		// ID:          project.ProjectID,
	// 		Name:        project.Name,
	// 		Description: project.Description,
	// 		Status:      project.Status,
	// 		ContractURL: project.ContractURL,
	// 		StartDate:   project.StartDate,
	// 		EndDate:     project.EndDate,
	// 		// QuotationID: responses.NullableUUID{UUID: project.QuotationID, Valid: project.QuotationID != uuid.Nil},
	// 		// ContractID:  responses.NullableUUID{UUID: project.ContractID, Valid: project.ContractID != uuid.Nil},
	// 		// InvoiceID:   responses.NullableUUID{UUID: project.InvoiceID, Valid: project.InvoiceID != uuid.Nil},
	// 		// BID:         responses.NullableUUID{UUID: project.BID, Valid: project.BID != uuid.Nil},
	// 		// ClientID:    responses.NullableUUID{UUID: project.ClientID, Valid: project.ClientID != uuid.Nil},
	// 		CreatedAt: project.CreatedAt,
	// 		UpdatedAt: project.UpdatedAt,
	// 	})
	// }

	return c.JSON(fiber.Map{
		"data": projects,
	})
}

func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	project, err := h.projectUsecase.GetProject(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	return c.JSON(fiber.Map{
		"data": responses.ProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			ContractURL: project.ContractURL,
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
			QuotationID: responses.NullableUUID{UUID: project.QuotationID, Valid: project.QuotationID != uuid.Nil},
			ContractID:  responses.NullableUUID{UUID: project.ContractID, Valid: project.ContractID != uuid.Nil},
			InvoiceID:   responses.NullableUUID{UUID: project.InvoiceID, Valid: project.InvoiceID != uuid.Nil},
			BID:         responses.NullableUUID{UUID: project.BID, Valid: project.BID != uuid.Nil},
			ClientID:    responses.NullableUUID{UUID: project.ClientID, Valid: project.ClientID != uuid.Nil},
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		},
	})
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var req requests.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	project, err := h.projectUsecase.UpdateProject(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update project"})
	}

	return c.JSON(fiber.Map{
		"message": "Project updated successfully",
		"data": responses.UpdateProjectResponse{
			ID:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			Status:      project.Status,
			ContractURL: project.ContractURL,
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
		},
	})
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	err = h.projectUsecase.DeleteProject(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete project"})
	}

	return c.JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}
