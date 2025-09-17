package db_connect

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/pdcgo/shared/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewProductionDatabase(appname string, cfg *configs.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        cfg.ToDsn(appname),
	}), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, err
}
