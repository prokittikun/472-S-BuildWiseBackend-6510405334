package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(context.Context, requests.LoginRequest) (*responses.UserResponse, string, error)
	Register(context.Context, requests.RegisterRequest) error
}

type userUsecase struct {
	userRepo    repositories.UserRepository
	jwtSecret   []byte
	jwtDuration time.Duration
}

func NewUserUsecase(userRepo repositories.UserRepository, jwtSecret string, jwtDuration time.Duration) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
		jwtDuration: jwtDuration,
	}
}

func (uu *userUsecase) generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"exp":      time.Now().Add(uu.jwtDuration).Unix(),
	})

	return token.SignedString(uu.jwtSecret)
}

func (uu *userUsecase) Login(
	ctx context.Context,
	loginRequest requests.LoginRequest,
) (*responses.UserResponse, string, error) {
	user, err := uu.userRepo.GetByUsername(ctx, loginRequest.Username)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := uu.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	userResponse := &responses.UserResponse{
		ID:        user.UserID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Tel:       user.Tel,
	}
	return userResponse, token, nil
}

func (uu *userUsecase) Register(
	ctx context.Context,
	registerRequest requests.RegisterRequest,
) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	registerRequest.Password = string(hashedPassword)
	return uu.userRepo.CreateUser(ctx, registerRequest)
}
