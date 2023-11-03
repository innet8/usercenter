package migrations

import (
	"authentik-go/app/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func (s InitDatabase) AddTableAuthentikCoreUser() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023110200",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.AuthentikCoreUser{})
		},
	}
}
