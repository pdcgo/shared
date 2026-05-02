package moretest_mock

import (
	"errors"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/bigquery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbScenario func(t *testing.T, handler func(tx *gorm.DB))

func MockPostgresDatabase(scenario *DbScenario) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		tempDb, err := db_connect.ConnectLocalDatabaseTest()
		if err != nil {
			log.Fatalf("database: failed to connect: %v", err)
		}

		*scenario = func(t *testing.T, handler func(tx *gorm.DB)) {
			t.Helper()
			tempDb.Transaction(func(tx *gorm.DB) error {
				handler(tx)
				return errors.New("rollback") // for keep database test clean up
			})
		}

		return nil
	}
}

func MockSqliteDatabase(db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {

		id := uuid.New()
		os.MkdirAll("/tmp/db_test", os.ModeDir)
		fname := "/tmp/db_test/" + id.String() + ".db"

		dbs, err := gorm.Open(sqlite.Open(fname), &gorm.Config{})

		*db = *dbs

		assert.Nil(t, err)

		return func() error {
			dbInstance, _ := db.DB()
			_ = dbInstance.Close()
			os.Remove(fname)
			return nil
		}
	}
}

var bqdb *gorm.DB
var initbqdb sync.Once

func MockBigqueryDatabase(cdn string, db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		initbqdb.Do(func() {
			var err error
			err = moretest.CheckGCPAuth()
			if err != nil {
				return
			}

			bqdb, err = gorm.Open(bigquery.Open(cdn), &gorm.Config{})
			if err != nil {
				panic(err)
			}
		})

		if bqdb == nil {
			return nil
		}

		*db = *bqdb
		return nil
	}
}

func DisableForeignKey(db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		err := db.Exec("PRAGMA foreign_keys = OFF").Error
		assert.Nil(t, err)

		return func() error {
			err := db.Exec("PRAGMA foreign_keys = ON").Error
			assert.Nil(t, err)
			return nil
		}
	}
}
