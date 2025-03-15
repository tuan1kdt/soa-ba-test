package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	StatusUnknown    ProductStatus = "Unknown"
	StatusAvailable  ProductStatus = "Available"
	StatusOnOrDer    ProductStatus = "On Order"
	StatusOutOfStock ProductStatus = "Out of Stock"
)

func (s ProductStatus) String() string {
	return string(s)
}

func ParseProductStatus(status string) ProductStatus {
	switch status {
	case "Available":
		return StatusAvailable
	case "On Order":
		return StatusOnOrDer
	case "Out of Stock":
		return StatusOutOfStock
	default:
		return StatusUnknown
	}
}

// Product is an entity that represents a product
type Product struct {
	ID         uuid.UUID
	Reference  string
	Name       string
	AddedDate  time.Time
	Status     ProductStatus
	CategoryID *uuid.UUID
	Price      float64
	StockCity  string
	SupplierID *uuid.UUID
	Quantity   int

	Category *Category
}
