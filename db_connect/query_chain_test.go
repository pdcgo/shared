package db_connect_test

import (
	"log"
	"testing"

	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type DataTemp struct {
	gorm.Model
	Name string
}

func TestChainQuery(t *testing.T) {
	var db gorm.DB
	moretest.Suite(t, "test Chain", moretest.SetupListFunc{
		moretest_mock.MockSqliteDatabase(&db),
		func(t *testing.T) func() error {
			err := db.AutoMigrate(
				&DataTemp{},
			)

			assert.Nil(t, err)
			return nil
		},
	},
		func(t *testing.T) {
			_, err := db_connect.NewQueryChain(&db,
				func(db *gorm.DB, next db_connect.NextFunc) db_connect.NextFunc {
					return func(query *gorm.DB) (*gorm.DB, error) {
						log.Println("chain 1")

						return next(query)
					}
				},
				func(db *gorm.DB, next db_connect.NextFunc) db_connect.NextFunc {
					return func(query *gorm.DB) (*gorm.DB, error) {
						_, err := next(query)
						log.Println("chain 2", err)

						return next(query)
					}
				},
				func(db *gorm.DB, next db_connect.NextFunc) db_connect.NextFunc {
					return func(query *gorm.DB) (*gorm.DB, error) {
						log.Println("chain 3")
						return next(query)
					}
				},
			)

			assert.Nil(t, err)
		},
	)

}
