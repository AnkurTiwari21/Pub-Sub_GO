package models

import (
	"time"

	"gorm.io/gorm"
)

type Queue struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Exchanges []Exchange `gorm:"many2many:bindings;"` 
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
