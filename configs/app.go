package configs

import (
	"fmt"

	"github.com/pdcgo/shared/interfaces/wallet_iface"
	"github.com/pdcgo/shared/pkg/secret"
)

type AppConfig struct {
	JwtSecret string         `yaml:"jwt_secret"`
	Database  DatabaseConfig `yaml:"database"`

	StatService       StatService                      `yaml:"stat_service"`
	WalletService     wallet_iface.WalletServiceConfig `yaml:"wallet_service"`
	TrackService      TrackService                     `yaml:"track_service"`
	AccountingService AccountingService                `yaml:"accounting_service"`
	WithdrawalService WithdrawalService                `yaml:"withdrawal_service"`
	CacheService      CacheService                     `yaml:"cache_service"`

	DispatcherConfig DispatcherConfig `yaml:"dispatcher_config"`
}

type CacheService struct {
	Endpoint string `yaml:"endpoint"`
}

type DispatcherQueue string

const (
	SlowQueue DispatcherQueue = "slow-event-task"
)

type DispatcherConfig struct {
	ProjectID string `yaml:"project_id"`
	Location  string `yaml:"location"`
}

func (c *DispatcherConfig) GetPath(key DispatcherQueue) string {
	return fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.ProjectID, c.Location, key)
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

type WithdrawalService struct {
	Endpoint string `yaml:"endpoint"`
}

type AccountingService struct {
	Endpoint string `yaml:"endpoint"`
}

type StatService struct {
	SubID string `yaml:"sub_id"`
}

type TrackService struct {
	SubID string `yaml:"sub_id"`
}

func NewProductionConfig() (*AppConfig, error) {
	var cfg AppConfig
	var sec *secret.Secret
	var err error

	// getting configuration
	sec, err = secret.GetSecret("app_config_prod", "latest")
	if err != nil {
		panic(err)
	}

	err = sec.YamlDecode(&cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
