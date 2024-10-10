package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := ur.db.GetContext(ctx, user, `SELECT * FROM  "User" WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, req requests.RegisterRequest) error {
	user := &models.User{
		UserID:    uuid.New(),
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Tel:       req.Tel,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO "User" (user_id, username, password, first_name, last_name, email, tel, created_at, updated_at)
			  VALUES (:user_id, :username, :password, :first_name, :last_name, :email, :tel, :created_at, :updated_at)`

	_, err := ur.db.NamedExecContext(ctx, query,
		map[string]interface{}{
			"user_id":    user.UserID,
			"username":   user.Username,
			"password":   user.Password,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"tel":        user.Tel,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
