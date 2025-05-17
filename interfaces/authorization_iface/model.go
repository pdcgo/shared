package authorization_iface

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/gorm"
)

type ExpiredToken struct {
	LastCreated int64
	LastReset   int64
}

type Identity interface {
	GetToken(appname, secret string) (string, error)
	GetAgentType() identity_iface.AgentType
	IsSuperUser() bool
	IdentityID() uint
	GetUserID() uint
	IsTokenExpired(tx *gorm.DB) (bool, error)
	HasRole(tx *gorm.DB, domainID uint, keyname string) (bool, error)
	GetExpired(tx *gorm.DB) (*ExpiredToken, error)
}

type Entity interface {
	GetEntityID() string
}

type Policy int

func (p Policy) ToBool() bool {
	return p == 1
}

const (
	Allow Policy = 1
	Deny  Policy = 0
)

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Read   Action = "read"
	Delete Action = "delete"
)

type CacheAuthorize struct {
	ID     string `gorm:"primarykey,autoIncrement:false"`
	UserID uint

	Policy Policy
}

type Role struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Key      string `json:"key" gorm:"index:domain_key,unique"`
	DomainID uint   `json:"domain_id" gorm:"index:domain_key,unique"`

	Permissions []*Permission     `json:"permission" gorm:"foreignKey:RoleID"`
	Caches      []*CacheAuthorize `gorm:"many2many:role_caches;" json:"-"`

	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func (r *Role) GetEntityID() string {
	return "role"
}

type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	RoleID uint `gorm:"index:user_role,unique"`
	UserID uint `gorm:"index:user_role,unique"`

	Role *Role
}

func (u *UserRole) GetEntityID() string {
	return "user_role"
}

type Permission struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RoleID   uint   `gorm:"index:permission_unique,unique" json:"role_id"`
	DomainID uint   `gorm:"index:permission_unique,unique" json:"domain_id"`
	EntityID string `gorm:"index:permission_unique,unique" json:"entity_id"`
	Action   Action `gorm:"index:permission_unique,unique" json:"action"`
	Policy   Policy `json:"policy"`

	Role *Role `json:"-"`
}

func (p *Permission) GetEntityID() string {
	return "permission"
}

type CheckPermission struct {
	DomainID uint     `gorm:"primaryKey" json:"domain_id"`
	Actions  []Action `gorm:"primaryKey" json:"action"`
}

type CheckPermissionGroup map[Entity]*CheckPermission

func (check CheckPermissionGroup) Permission() []*Permission {
	hasil := []*Permission{}

	for ent, value := range check {
		for _, action := range value.Actions {
			hasil = append(hasil, &Permission{
				DomainID: value.DomainID,
				EntityID: ent.GetEntityID(),
				Policy:   Allow,
				Action:   action,
			})
		}
	}

	return hasil
}

func (check CheckPermissionGroup) GetDomainIDs() []uint {
	hasil := []uint{}

	for _, value := range check {
		hasil = append(hasil, value.DomainID)
	}

	return hasil
}

func (check CheckPermissionGroup) GetEntityIDs() []string {
	hasil := []string{}
	for ent := range check {
		hasil = append(hasil, ent.GetEntityID())
	}

	return hasil
}

func (check CheckPermissionGroup) GetActions() []Action {
	hasil := []Action{}

	actionmap := map[Action]bool{}

	for _, value := range check {
		for _, action := range value.Actions {
			actionmap[action] = true
		}
	}

	for key := range actionmap {
		hasil = append(hasil, key)
	}

	return hasil
}

func (li CheckPermissionGroup) GetID() (string, error) {
	data, err := json.Marshal(li.Permission())
	if err != nil {
		return "", err
	}
	idstr := fmt.Sprintf("%x", md5.Sum(data))
	index := len(idstr) - 1
	return idstr[index-8 : index], nil
}

type PermissionQuery interface {
	Permission(agent identity_iface.Agent) CheckPermissionGroup
}

type PermissionError struct {
	Err              error
	NeedPermissions  []*Permission `json:"need_permission"`
	ActualPermission []*Permission `json:"actual_permission"`
}

func (permerr *PermissionError) Error() string {
	if len(permerr.ActualPermission) == 0 {
		return "permission kosong"
	}

	return "permission restricted"
}

func (err *PermissionError) Unwrap() error {
	return err.Err
}
