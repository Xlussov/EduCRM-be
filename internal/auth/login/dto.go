package login

import "github.com/google/uuid"

type Request struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

type UserDTO struct {
	ID       uuid.UUID   `json:"id"`
	Role     string      `json:"role"`
	Branches []uuid.UUID `json:"branches"`
}
