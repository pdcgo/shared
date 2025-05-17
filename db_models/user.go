package db_models

import (
	"errors"
	"time"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/gorm"
)

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

// GetToken implements authorization_iface.Identity.
func (u *User) GetToken(appname string, secret string) (string, error) {
	return "", errors.New("cannot create token in plain db model")
}

type AppKey struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	AppName string `json:"app_name"`
	UserID  uint   `json:"user_id"`
	Key     string `json:"-"`

	User *User
}

// GetAgentType implements authorization_iface.Identity.
func (u *User) GetAgentType() identity_iface.AgentType {
	return identity_iface.SystemAgent
}

// IsTokenExpired implements authorization_iface.Identity.
func (u *User) IsTokenExpired(tx *gorm.DB) (bool, error) {
	return false, nil
}

// GetEntityID implements authorization.Entity.
func (u *User) GetEntityID() string {
	return "user"
}

// GetExpired implements authorization_iface.Identity.
func (u *User) GetExpired(tx *gorm.DB) (*authorization_iface.ExpiredToken, error) {
	return &authorization_iface.ExpiredToken{
		LastCreated: u.LastCreated,
		LastReset:   u.LastReset,
	}, nil
}

// GetUserID implements authorization_iface.Identity.
func (u *User) GetUserID() uint {
	return u.ID
}

// IdentityID implements authorization_iface.Identity.
func (u *User) IdentityID() uint {
	return u.ID
}

// IsSuperUser implements authorization_iface.Identity.
func (u *User) IsSuperUser() bool {
	return u.IsRoot
}

func (u *User) HasRole(db *gorm.DB, domainID uint, keyname string) (bool, error) {
	var roleCount int64

	roles := db.Model(&authorization_iface.Role{}).
		Select("id").
		Where("domain_id = ?", domainID).
		Where("key = ?", keyname)

	err := db.Model(&authorization_iface.UserRole{}).
		Where("role_id = (?)", roles).
		Where("user_id = ?", u.IdentityID()).
		Count(&roleCount).Error
	if err != nil {
		return false, err
	}

	hasRole := roleCount > 0
	return hasRole, err
}

type UserTeam struct {
	ID     uint   `gorm:"primarykey" json:"id"`
	UserID uint   `json:"user_id" gorm:"index:userteam_uniquea,unique"`
	TeamID uint   `json:"team_id" gorm:"index:userteam_uniquea,unique"`
	Alias  string `json:"alias" gorm:"index:userteam_uniquea,unique"`

	User *User `json:"-"`
	Team *Team `json:"team"`
}
