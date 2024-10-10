package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"boonkosang/internal/usecase"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JobHandler struct {
	jobUsecase usecase.JobUsecase
}

func NewJobHandler(jobUsecase usecase.JobUsecase) *JobHandler {
	return &JobHandler{
		jobUsecase: jobUsecase,
	}
}

func (h *JobHandler) JobRoutes(app *fiber.App) {
	app.Post("/jobs", h.CreateJob)
	app.Get("/jobs", h.ListJobs)
	app.Get("/jobs/:id", h.GetJob)
	app.Put("/jobs/:id", h.UpdateJob)
	app.Delete("/jobs/:id", h.DeleteJob)
}

func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	var req requests.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	job, err := h.jobUsecase.CreateJob(c.Context(), req)
	fmt.Println(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create job"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Job created successfully",
		"data":    responses.NewJobResponse(job),
	})
}

func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job ID"})
	}

	job, err := h.jobUsecase.GetJob(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}

	return c.JSON(fiber.Map{
		"data": responses.NewJobResponse(job),
	})
}

func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	jobs, err := h.jobUsecase.ListJobs(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch jobs"})
	}

	jobResponses := make([]responses.JobResponse, len(jobs))
	for i, job := range jobs {
		jobResponses[i] = responses.NewJobResponse(job)
	}

	return c.JSON(fiber.Map{
		"data": jobResponses,
	})
}

func (h *JobHandler) UpdateJob(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job ID"})
	}

	var req requests.UpdateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	job, err := h.jobUsecase.UpdateJob(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update job"})
	}

	return c.JSON(fiber.Map{
		"message": "Job updated successfully",
		"data":    responses.NewJobResponse(job),
	})
}

func (h *JobHandler) DeleteJob(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job ID"})
	}

	err = h.jobUsecase.DeleteJob(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete job"})
	}

	return c.JSON(fiber.Map{
		"message": "Job deleted successfully",
	})
}
