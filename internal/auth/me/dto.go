package me

import (
	"github.com/Xlussov/EduCRM-be/internal/domain"
	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID   `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Phone     string      `json:"phone"`
	Role      domain.Role `json:"role"`
}
