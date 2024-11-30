package models

import (
	"time"

	"gorm.io/gorm"
)

type Exchange struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Queues    []Queue `gorm:"many2many:bindings;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
