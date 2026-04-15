package domain

import (
	"fmt"
	"strings"
)

func ParseEntityStatus(raw string) (*EntityStatus, error) {
	if raw == "" {
		return nil, nil
	}

	status := EntityStatus(strings.ToUpper(raw))
	if status != StatusActive && status != StatusArchived {
		return nil, fmt.Errorf("%w: status must be ACTIVE or ARCHIVED", ErrInvalidInput)
	}

	return &status, nil
}
