package responses

import "github.com/google/uuid"

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Tel       string    `json:"tel"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	user        UserResponse `json:"user"`
}
