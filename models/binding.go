package models

import (
	"time"

	"gorm.io/gorm"
)

type Binding struct {
	ExchangeId uint `gorm:"primaryKey"`
	QueueId    uint `gorm:"primaryKey"`

	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
