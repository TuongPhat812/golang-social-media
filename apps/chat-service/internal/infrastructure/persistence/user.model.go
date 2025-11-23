package persistence

import (
	"time"
)

type UserModel struct {
	ID        string    `gorm:"column:id;type:text;primaryKey"`
	Email     string    `gorm:"column:email;type:text;not null"`
	Name      string    `gorm:"column:name;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (UserModel) TableName() string {
	return "users"
}

