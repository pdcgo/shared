package moretest_mock

import (
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/bigquery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MockPostgresDatabase(db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {
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
