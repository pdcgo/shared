package db_connect

import (
	"context"
	"fmt"
	"reflect"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/pkg/gorm_commenter"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

func NewRedisDatabase(cfg *configs.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

type TimestampSerializer struct{}

func init() {
	schema.RegisterSerializer("timestamptz", TimestampSerializer{})
}

func (TimestampSerializer) Scan(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	dbValue interface{},
) error {
	if dbValue == nil {
		field.ReflectValueOf(ctx, dst).Set(reflect.Zero(field.FieldType))
		return nil
	}

	var t time.Time

	switch v := dbValue.(type) {
	case time.Time:
		t = v
	case string:
		parsed, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return fmt.Errorf("timestamptz: failed to parse string %q: %w", v, err)
		}
		t = parsed
	case []byte:
		parsed, err := time.Parse(time.RFC3339Nano, string(v))
		if err != nil {
			return fmt.Errorf("timestamptz: failed to parse bytes: %w", err)
		}
		t = parsed
	default:
		return fmt.Errorf("timestamptz: unsupported type %T", dbValue)
	}

	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(timestamppb.New(t)))
	return nil
}

// ✅ return (interface{}, error) not (driver.Value, error)
func (TimestampSerializer) Value(
	ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	fieldValue interface{},
) (interface{}, error) {
	if fieldValue == nil {
		return nil, nil
	}

	ts, ok := fieldValue.(*timestamppb.Timestamp)
	if !ok {
		return nil, fmt.Errorf("timestamptz: expected *timestamppb.Timestamp, got %T", fieldValue)
	}

	if ts == nil {
		return nil, nil
	}

	if err := ts.CheckValid(); err != nil {
		return nil, fmt.Errorf("timestamptz: invalid timestamp: %w", err)
	}

	return ts.AsTime(), nil // return time.Time, GORM handles the rest
}
