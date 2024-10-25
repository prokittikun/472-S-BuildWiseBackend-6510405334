package rest

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) UserRoutes(app *fiber.App) {
	app.Post("/login", h.Login)
	app.Post("/register", h.Register)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest requests.LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	loginResponse, err := uh.userUsecase.Login(c.Context(), loginRequest)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.JSON(
		fiber.Map{
			"message": "Login successful",
			"data":    loginResponse,
		},
	)
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	var registerRequest requests.RegisterRequest

	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := uh.userUsecase.Register(c.Context(), registerRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}
