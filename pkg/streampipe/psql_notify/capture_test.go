package psql_notify_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/streampipe/psql_notify"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LogStock struct {
	ID   uint `gorm:"primarykey" json:"id"`
	Data string
}

func TestCapture(t *testing.T) {
	var db gorm.DB

	var migration moretest.SetupFunc = func(t *testing.T) func() error {
		err := db.AutoMigrate(
			&LogStock{},
		)
		assert.Nil(t, err)
		return nil
	}

	moretest.Suite(t, "testing capture",
		moretest.SetupListFunc{
			ConnectTestPostgres(&db),
			migration,
		},
		func(t *testing.T) {
			dsn := "host=localhost user=admintest password=admintest dbname=postgres port=5432 sslmode=disable"
			cap := psql_notify.NewCapturePostgres(psql_notify.DBTypePostgres, dsn)

			ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
			datas := make(chan string, 1)
			err := cap.
				PrepareChannel("log_stocks", "log_stocks_channel").
				Listen(ctx, "log_stocks_channel", true, func(raw string) {
					datas <- raw
					cancel()
				}, nil).
				Err()

			assert.Nil(t, err)

			data := LogStock{
				Data: "asdasdasd",
			}
			err = db.Save(&data).Error
			assert.Nil(t, err)

			raw := <-datas
			assert.NotEmpty(t, raw)
		},
	)

}

func ConnectTestPostgres(db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		DetectPostgres(t)

		dsn := "host=localhost user=admintest password=admintest dbname=postgres port=5432 sslmode=disable"
		mockdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			assert.Nil(t, err)
		}

		*db = *mockdb
		return nil
	}
}

func DetectPostgres(t *testing.T) {
	host := "localhost"
	port := 5432

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, time.Second*5)
	if err != nil {
		// t.Error(err)
		t.Skip("postgres testing environtment not detected")
		return
	}

	conn.Close()
}
