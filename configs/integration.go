package configs

type LegacyService struct {
	Endpoint string `yaml:"endpoint"`
}

func (s *LegacyService) JoinPath(path string) string {

	return s.Endpoint + path
}

type AppProductionTestConfig struct {
	JwtSecret         string            `yaml:"jwt_secret"`
	Database          DatabaseConfig    `yaml:"database"`
	LegacyService     LegacyService     `yaml:"legacy_service"`
	AccountingService AccountingService `yaml:"accounting_service"`
}
