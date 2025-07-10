package configs

import (
	"fmt"

	"github.com/pdcgo/shared/interfaces/wallet_iface"
)

type AppConfig struct {
	JwtSecret string         `yaml:"jwt_secret"`
	Database  DatabaseConfig `yaml:"database"`

	StatService   StatService                      `yaml:"stat_service"`
	WalletService wallet_iface.WalletServiceConfig `yaml:"wallet_service"`
	TrackService  TrackService                     `yaml:"track_service"`
}

type DatabaseConfig struct {
	DBName     string `yaml:"DB_NAME"`
	DBUser     string `yaml:"DB_USER"`
	DBPass     string `yaml:"DB_PASS"`
	DBInstance string `yaml:"DB_INSTANCE"`
}

func (cfg *DatabaseConfig) ToDsn(appName string) string {
	// dsn := "user=postgres password=postgres dbname=postgres host=postgres sslmode=disable"
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s application_name=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBInstance,
		appName,
	)
}

type StatService struct {
	SubID string `yaml:"sub_id"`
}

type TrackService struct {
	SubID string `yaml:"sub_id"`
}
