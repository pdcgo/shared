package db_connect

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/pkg/secret"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBCLuster interface {
	Initialize() error
	Read(ctx context.Context) *gorm.DB
	Write(ctx context.Context) *gorm.DB
}

type CloudSqlInstance struct {
	ProjectID string
	Location  string
	Name      string
}

func (i *CloudSqlInstance) String() string {
	return fmt.Sprintf("%s:%s:%s", i.ProjectID, i.Location, i.Name)
}

func ParseCloudSqlInstance(instance string) *CloudSqlInstance {
	datas := strings.Split(instance, ":")
	return &CloudSqlInstance{
		ProjectID: datas[0],
		Location:  datas[1],
		Name:      datas[2],
	}
}

type DBModifier func(db *gorm.DB) *gorm.DB

type prodCluster struct {
	ctx          context.Context
	projectID    string
	appname      string
	srv          *sqladmin.Service
	cfg          *configs.DatabaseConfig
	modifier     DBModifier
	next         int
	replicaCount int
	primary      *gorm.DB
	replicas     []*gorm.DB
}

func (p *prodCluster) connect(cfg *configs.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        cfg.ToDsn(p.appname),
	}), &gorm.Config{})
	if err != nil {
		return db, err
	}

	if p.modifier != nil {
		db = p.modifier(db)
	}

	return db, err

}

// Initialize implements DBCLuster.
func (p *prodCluster) Initialize() error {
	var err error
	// connecting primary
	slog.Info("connecting primary database")
	p.primary, err = p.connect(p.cfg)
	if err != nil {
		return err
	}

	// Get instance info
	instn := ParseCloudSqlInstance(p.cfg.DBInstance)
	inst, err := p.srv.Instances.Get(p.projectID, instn.Name).Context(p.ctx).Do()
	if err != nil {
		return err
	}

	p.replicaCount = len(inst.ReplicaNames)
	for _, r := range inst.ReplicaNames {

		iname := CloudSqlInstance{
			ProjectID: p.projectID,
			Location:  instn.Location,
			Name:      r,
		}
		slog.Info("connecting replica database", slog.String("replica_name", r))
		db, err := p.connect(&configs.DatabaseConfig{
			DBName:     p.cfg.DBName,
			DBUser:     p.cfg.DBUser,
			DBPass:     p.cfg.DBPass,
			DBInstance: iname.String(),
		})

		if err != nil {
			return err
		}

		p.replicas = append(p.replicas, db)
	}

	return nil
}

// Read implements DBCLuster.
func (p *prodCluster) Read(ctx context.Context) *gorm.DB {
	if p.replicaCount == 0 {
		return p.primary
	}

	db := p.replicas[p.next]
	p.next = (p.next + 1) % p.replicaCount

	return db.WithContext(ctx)
}

// Write implements DBCLuster.
func (p *prodCluster) Write(ctx context.Context) *gorm.DB {
	return p.primary.WithContext(ctx)
}

func NewProdCluster(appname string, cfg *configs.DatabaseConfig, modifier DBModifier) DBCLuster {
	ctx := context.Background()
	sqlService, err := sqladmin.NewService(ctx, option.WithScopes(sqladmin.CloudPlatformScope))
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		var prodCfg configs.AppConfig
		var sec *secret.Secret

		// getting configuration
		sec, err = secret.GetSecret("app_config_prod", "latest")
		if err != nil {
			panic(err)
		}

		err = sec.YamlDecode(&prodCfg)
		if err != nil {
			panic(err)
		}

		cfg = &prodCfg.Database
	}

	return &prodCluster{
		projectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		srv:       sqlService,
		appname:   appname,
		cfg:       cfg,
		replicas:  []*gorm.DB{},
		modifier:  modifier,
	}
}
