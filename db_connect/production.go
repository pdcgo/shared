package db_connect

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/pkg/gorm_commenter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewProductionDatabase(appname string, cfg *configs.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        cfg.ToDsn(appname),
	}), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return db, err
	}

	err = db.Use(gorm_commenter.NewCommentClausePlugin())
	return db, err
}
