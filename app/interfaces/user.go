package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type UserReq struct {
	ID                 int                    `gorm:"primary_key" json:"id"`
	Password           string                 `gorm:"type:varchar(128);not null" json:"password"`
	LastLogin          time.Time              `gorm:"type:timestamptz(6)" json:"last_login"`
	Username           string                 `gorm:"type:varchar(150);not null" json:"username"`
	FirstName          string                 `gorm:"type:varchar(150);not null" json:"first_name"`
	LastName           string                 `gorm:"type:varchar(150);not null" json:"last_name"`
	Email              string                 `gorm:"type:varchar(254);not null" json:"email"`
	IsActive           bool                   `gorm:"not null" json:"is_active"`
	DateJoined         time.Time              `gorm:"type:timestamptz(6);not null" json:"date_joined"`
	UUID               uuid.UUID              `gorm:"not null" json:"uuid"`
	Name               string                 `gorm:"type:text;not null" json:"name"`
	PasswordChangeDate time.Time              `gorm:"type:timestamptz(6);not null" json:"password_change_date"`
	Attributes         map[string]interface{} `gorm:"type:jsonb;not null" json:"attributes"`
	Path               string                 `gorm:"type:text;not null" json:"path"`
	Type               string                 `gorm:"type:text;not null" json:"type"`
	Source             string                 `gorm:"type:varchar(50);default:'';comment:来源" json:"source"`
}

type UserLoginReq struct {
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

type UserRegReq struct {
	Email    string `json:"email"`    //邮箱
	Password string `json:"password"` //密码
	Source   string `json:"source"`   //来源
}
