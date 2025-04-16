package gorm_commenter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/shared/pkg/gorm_commenter"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMigrate(db *gorm.DB) moretest.SetupFunc {
	return func(t *testing.T) func() error {
		err := db.AutoMigrate(&ModUser{})
		assert.Nil(t, err)

		return func() error {
			return nil
		}
	}
}

type ModUser struct {
	*gorm.Model
	Name string
}

func TestGormTraceTag(t *testing.T) {
	var db gorm.DB

	moretest.Suite(t, "testing sqlcomenter",
		moretest.SetupListFunc{
			moretest_mock.MockSqliteDatabase(&db),
			setupMigrate(&db),
		},
		func(t *testing.T) {
			db.Use(gorm_commenter.NewCommentClausePlugin())
			dd := "action"
			ctx := context.WithValue(context.Background(), dd, "test")
			sql := db.WithContext(ctx).ToSQL(func(tx *gorm.DB) *gorm.DB {
				user := []ModUser{}
				return tx.Model(&ModUser{}).Find(&user)
			})
			assert.Contains(t, sql, "action=test")
			// t.Error(sql)
		})

	moretest.Suite(t, "testing dengan gin middleware",
		moretest.SetupListFunc{
			moretest_mock.MockSqliteDatabase(&db),
			setupMigrate(&db),
		},
		func(t *testing.T) {
			r := gin.Default()
			r.Use(gorm_commenter.GormCommenterMiddleware())
			db.Use(gorm_commenter.NewCommentClausePlugin())
			r.GET("/test", func(ctx *gin.Context) {
				sql := db.WithContext(ctx).ToSQL(func(tx *gorm.DB) *gorm.DB {
					user := []ModUser{}
					return tx.Model(&ModUser{}).Find(&user)
				})

				ctx.Data(http.StatusOK, "", []byte(sql))

			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			sql := w.Body.String()

			assert.Contains(t, sql, "route=/test")
			// t.Error(sql)

		})
}
