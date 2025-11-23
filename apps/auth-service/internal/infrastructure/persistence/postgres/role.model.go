package postgres

import (
	"time"
)

type RoleModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;type:text;not null;uniqueIndex"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (RoleModel) TableName() string {
	return "roles"
}

