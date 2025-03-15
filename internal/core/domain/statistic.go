package domain

import "github.com/google/uuid"

type StatisticSupplierProduct struct {
	SupplierID   uuid.UUID
	SupplierName string
	Percentage   float64
}

type StatisticCategoryProduct struct {
	CategoryID   uuid.UUID
	CategoryName string
	Percentage   float64
}
