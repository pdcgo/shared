package db_models

import "time"

type UserType string

const (
	UserCommon UserType = "common"
	UserSystem UserType = "system"
)

type User struct {
	ID             uint     `json:"id" gorm:"primarykey"`
	UserType       UserType `json:"user_type"`
	Name           string   `json:"name"`
	ProfilePicture string   `json:"profile_picture"`
	Username       string   `json:"username" gorm:"index:username_unique,unique"`
	Password       string   `json:"-"`
	Email          string   `json:"email" gorm:"index:email_unique,unique"`
	PhoneNumber    string   `json:"phone_number"`
	IsSuspended    bool     `json:"is_suspend"`
	IsRoot         bool     `json:"is_root"`

	LastCreated       int64     `json:"-"`
	LastReset         int64     `json:"-"`
	LastPasswordReset time.Time `json:"last_password_reset"`
	InvitationCode    string    `json:"-"`
	CreatedAt         time.Time `json:"created"`
}

func (User) TableName() string {
	return "users"
}

type AppKey struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	AppName string `json:"app_name"`
	UserID  uint   `json:"user_id"`
	Key     string `json:"-"`

	User *User
}
