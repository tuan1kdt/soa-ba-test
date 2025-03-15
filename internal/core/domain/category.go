package domain

import (
	"github.com/google/uuid"
)

// Category is an entity that represents a category of product
type Category struct {
	ID   uuid.UUID
	Name string
}
