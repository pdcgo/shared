package main

import (
	"context"
	"log"

	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/pkg/secret"
)

func main() {
	var cfg configs.AppConfig
	var sec *secret.Secret
	var err error

	// getting configuration
	sec, err = secret.GetSecret("app_config_prod", "latest")
	if err != nil {
		panic(err)
	}

	err = sec.YamlDecode(&cfg)
	if err != nil {
		panic(err)
	}

	cluster := db_connect.NewProdCluster("test", &cfg.Database, nil)
	err = cluster.Initialize()
	if err != nil {
		panic(err)
	}

	db := cluster.Read(context.Background())

	for range [20]int{} {
		teams := []*db_models.Team{}

		err := db.Model(&db_models.Team{}).Find(&teams).Error
		if err != nil {
			panic(err)
		}

		log.Println("requesting team")
	}

}
