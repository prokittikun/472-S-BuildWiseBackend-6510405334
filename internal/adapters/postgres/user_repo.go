package postgres

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	query := `SELECT * FROM "User" WHERE username = $1`
	err := ur.db.GetContext(ctx, user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
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
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Tel:       sql.NullString{String: req.Tel, Valid: req.Tel != ""},
	}

	query := `
        INSERT INTO "User" (
            user_id, username, password, first_name, last_name, 
            email, tel
        ) VALUES (
            :user_id, :username, :password, :first_name, :last_name, 
            :email, :tel
        )`

	_, err := ur.db.NamedExecContext(ctx, query, user)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return errors.New("username already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
