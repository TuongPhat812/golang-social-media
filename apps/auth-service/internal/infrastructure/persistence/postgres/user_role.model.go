package postgres

import (
	"time"
)

type UserRoleModel struct {
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID    string    `gorm:"column:role_id;type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (UserRoleModel) TableName() string {
	return "user_roles"
}

