package postgres

import (
	"time"
)

type PermissionModel struct {
	ID        string    `gorm:"column:id;type:uuid;primaryKey"`
	Name      string    `gorm:"column:name;type:text;not null"`
	Resource  string    `gorm:"column:resource;type:text;not null"`
	Action    string    `gorm:"column:action;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (PermissionModel) TableName() string {
	return "permissions"
}

