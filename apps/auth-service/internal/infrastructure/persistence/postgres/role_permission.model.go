package postgres

import (
	"time"
)

type RolePermissionModel struct {
	RoleID       string    `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionID string    `gorm:"column:permission_id;type:uuid;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

func (RolePermissionModel) TableName() string {
	return "role_permissions"
}

