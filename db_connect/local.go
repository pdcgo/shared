package db_connect

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectLocalDatabase() (*gorm.DB, error) {
	var err error
	host := getEnv("STAT_POSTGRES_HOST", "localhost")
	user := getEnv("STAT_POSTGRES_USER", "user")
	pass := getEnv("STAT_POSTGRES_PASSWORD", "password")
	dbname := getEnv("STAT_POSTGRES_DB", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta",
		host,
		user,
		pass,
		dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, err
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
