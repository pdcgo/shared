package db_connect_test

import (
	"testing"

	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/pkg/debugtool"
	"github.com/pdcgo/shared/pkg/moretest"
	"github.com/pdcgo/shared/pkg/moretest/moretest_mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type User struct {
	ID    uint
	Name  string
	Posts []*Post `gorm:"foreignkey:UserID"`
}

type Post struct {
	ID     uint
	UserID uint
	Title  string
}

func TestOneToMany(t *testing.T) {
	var db gorm.DB

	moretest.Suite(
		t,
		"testing preload",
		moretest.SetupListFunc{
			moretest_mock.MockSqliteDatabase(&db),
			func(t *testing.T) func() error {
				err := db.
					AutoMigrate(
						&User{},
						&Post{},
					)
				assert.Nil(t, err)

				users := []*User{
					{
						ID:   1,
						Name: "User 1",
					},
					{
						ID:   2,
						Name: "User 2",
					},
				}

				posts := []*Post{
					{
						ID:     1,
						UserID: 1,
						Title:  "Post 1",
					},
					{
						ID:     2,
						UserID: 2,
						Title:  "Post 2",
					},
				}

				for _, user := range users {
					err = db.Create(user).Error
					assert.Nil(t, err)
				}

				for _, post := range posts {
					err = db.Create(post).Error
					assert.Nil(t, err)
				}

				return nil
			},
		},
		func(t *testing.T) {

			var users []*User
			err := db.
				Find(&users).Error
			assert.Nil(t, err)

			err = db_connect.OneToMany(
				db.Debug(),
				users,
				&db_connect.Foreign{
					Field:  "Posts",
					RefKey: "ID",
					Fkey:   "UserID",
					ReferenceValueFunc: func(i int) any {
						return users[i].ID
					},
					QueryOpt: func(q *gorm.DB) *gorm.DB {
						return q.
							Table("posts as p")
					},
				},
			)
			assert.Nil(t, err)

			debugtool.LogJson(users)

			for _, user := range users {
				assert.NotNil(t, user.Posts)
				assert.Len(t, user.Posts, 1)
			}

		})

}
