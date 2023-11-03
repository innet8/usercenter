package model

import (
	"time"

	"github.com/google/uuid"
)

// CoreUser表
type AuthentikCoreUser struct {
	ID                 int       `gorm:"primary_key" json:"id"`
	Password           string    `gorm:"type:varchar(128);not null" json:"password,omitempty"`
	LastLogin          time.Time `gorm:"type:timestamptz(6)" json:"last_login"`
	Username           string    `gorm:"type:varchar(150);not null" json:"username"`
	FirstName          string    `gorm:"type:varchar(150);not null" json:"first_name"`
	LastName           string    `gorm:"type:varchar(150);not null" json:"last_name"`
	Email              string    `gorm:"type:varchar(254);not null" json:"email"`
	IsActive           bool      `gorm:"not null" json:"is_active"`
	DateJoined         time.Time `gorm:"type:timestamptz(6);not null" json:"date_joined"`
	UUID               uuid.UUID `gorm:"not null" json:"uuid"`
	Name               string    `gorm:"type:text;not null" json:"name"`
	PasswordChangeDate time.Time `gorm:"type:timestamptz(6);not null" json:"password_change_date"`
	Attributes         string    `gorm:"type:jsonb;not null" json:"attributes"`
	Path               string    `gorm:"type:text;not null" json:"path"`
	Type               string    `gorm:"type:text;not null" json:"type"`
	Source             string    `gorm:"type:varchar(50);default:'';comment:来源" json:"source"`
}
