package configs

import "fmt"

type AppConfig struct {
	Database    DatabaseConfig `yaml:"database"`
	StatService StatService    `yaml:"stat_service"`
}

type DatabaseConfig struct {
	DBName     string `yaml:"DB_NAME"`
	DBUser     string `yaml:"DB_USER"`
	DBPass     string `yaml:"DB_PASS"`
	DBInstance string `yaml:"DB_INSTANCE"`
}

func (cfg *DatabaseConfig) ToDsn() string {
	// dsn := "user=postgres password=postgres dbname=postgres host=postgres sslmode=disable"
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBInstance)
}

type StatService struct {
	SubID string `yaml:"sub_id"`
}
