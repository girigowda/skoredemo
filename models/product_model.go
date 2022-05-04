package models

import (
	"time"
)

type Product struct {
	// gorm.Model
	Code        string
	Price       uint
	PublishedAt time.Time
}

// Database TableName of this model
func (e *Example) ProductName() string {
	return "examples"
}
