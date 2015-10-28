package ksana

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `sql:"not null"`
	UpdatedAt time.Time `sql:"not null"`
}
