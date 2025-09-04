package authorization

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/pkg/ware_cache"
	"gorm.io/gorm"
)

var RootDomain uint = 1

type Identity interface {
	IsSuperUser() bool
	IdentityID() uint
	GetUserID() uint
	GetExpired(tx *gorm.DB) (*authorization_iface.ExpiredToken, error)
	HasRole(tx *gorm.DB, domainID uint, keyname string) (bool, error)
}

type CheckPermission struct {
	DomainID uint                         `gorm:"primaryKey" json:"domain_id"`
	Actions  []authorization_iface.Action `gorm:"primaryKey" json:"action"`
}

type Entity interface {
	GetEntityID() string
}

type CacheAuthorize struct {
	ID     string `gorm:"primarykey,autoIncrement:false" json:"id"`
	UserID uint   `json:"user_id"`

	Policy authorization_iface.Policy `json:"policy"`
}

type RoleCache struct {
	ID               uint `gorm:"primarykey"`
	CacheAuthorizeID string
	RoleID           uint

	CacheAuthorize *CacheAuthorize
	Role           *authorization_iface.Role
}

type Authorization struct {
	cache      ware_cache.Cache
	tx         *gorm.DB
	passphrase string
}

// AuthIdentityFromHeader implements authorization_iface.Authorization.
func (auth *Authorization) AuthIdentityFromHeader(header http.Header) authorization_iface.AuthIdentity {
	return NewAuthIdentityHttpHeader(auth, header, auth.passphrase)
}

// HasPermission implements authorization_iface.Authorization.
func (auth *Authorization) HasPermission(identity authorization_iface.Identity, perms authorization_iface.CheckPermissionGroup) error {
	_, err := auth.CheckPermission(identity, perms)
	return err
}

func (auth *Authorization) FlushCache() error {
	return auth.cache.Flush(context.Background())
}

func (auth *Authorization) GetPassPhrase() string {
	return auth.passphrase
}

func (auth *Authorization) RefreshCacheUser(userID uint) error {
	err := auth.tx.Transaction(func(tx *gorm.DB) error {

		subquery := tx.Model(&CacheAuthorize{}).
			Select("id").
			Where("user_id = ?", userID)
		err := tx.Model(&RoleCache{}).Where("cache_authorize_id IN (?)", subquery).Delete(&RoleCache{}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&CacheAuthorize{}).
			Where("user_id = ?", userID).
			Delete(&CacheAuthorize{}).
			Error

		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (auth *Authorization) refreshCacheRole(roleID uint) error {
	var err error

	err = auth.tx.Model(&RoleCache{}).Where("role_id = ?", roleID).Delete(&RoleCache{}).Error
	if err != nil {
		return err
	}
	role := authorization_iface.Role{}
	err = auth.tx.Model(&authorization_iface.Role{}).First(&role, roleID).Error
	if err != nil {
		return err
	}
	return nil
}

func (auth *Authorization) Domain(domainID uint) *Domain {
	return &Domain{
		auth: auth,
		ID:   domainID,
	}
}

var ErrPermission = errors.New("permission error")

type ListPermissions []*authorization_iface.Permission

func (li ListPermissions) GetRole() []*authorization_iface.Role {
	hasil := []*authorization_iface.Role{}

	maprole := map[uint]*authorization_iface.Role{}
	for _, item := range li {
		dd := item
		maprole[item.RoleID] = dd.Role
	}

	for _, d := range maprole {
		i := d
		hasil = append(hasil, i)
	}

	return hasil
}

func (auth *Authorization) ApiQueryCheckPermission(identity authorization_iface.Identity, query authorization_iface.PermissionQuery) (bool, error) {
	return auth.CheckPermission(identity, query.Permission(identity))
}

func (auth *Authorization) CheckPermission(identity authorization_iface.Identity, needPerms authorization_iface.CheckPermissionGroup) (bool, error) {
	var iscache bool = false

	if identity.IsSuperUser() {
		return false, nil
	}

	hasil := authorization_iface.PermissionError{
		Err:              nil,
		NeedPermissions:  needPerms.Permission(),
		ActualPermission: []*authorization_iface.Permission{},
	}

	cache := CacheAuthorize{}
	cacheID, err := needPerms.GetID(identity.GetUserID())
	if err != nil {
		hasil.Err = err
		return iscache, &hasil
	}

	err = auth.cache.Get(context.Background(), cacheID, &cache)

	// err = auth.tx.Model(&CacheAuthorize{}).
	// 	Where("id = ?", cacheID).
	// 	Where("user_id = ?", identity.GetUserID()).
	// 	First(&cache).Error

	if err == nil {
		iscache = true
		if cache.Policy == authorization_iface.Allow {
			return iscache, nil
		}
		hasil.Err = ErrPermission
		return iscache, &hasil
	}

	// if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	hasil.Err = ErrPermission
	// 	return iscache, &hasil
	// }

	userrole := auth.tx.Table("user_roles").Select("role_id").Where("user_id = ?", identity.IdentityID())

	query := auth.tx.Model(&authorization_iface.Permission{}).
		Preload("Role").
		Where("role_id IN (?)", userrole).
		Where("domain_id IN ?", needPerms.GetDomainIDs()).
		Where("entity_id IN ?", needPerms.GetEntityIDs()).
		Where("action IN ?", needPerms.GetActions())

	actualpermissions := ListPermissions{}
	err = query.Find(&actualpermissions).Error

	if err != nil {
		hasil.Err = err
		return iscache, &hasil
	}

	// checking expired
	// expired, err := identity.GetExpired(auth.tx)
	// if err != nil {
	// 	hasil.Err = err
	// 	return iscache, &hasil
	// }

	// if expired.LastReset > expired.LastCreated {
	// 	hasil.Err = errors.New("token expired")
	// 	return iscache, &hasil
	// }

	// creating cache
	cache = CacheAuthorize{
		ID:     cacheID,
		UserID: identity.IdentityID(),
		Policy: authorization_iface.Allow,
	}

	if len(needPerms) != len(actualpermissions) {
		hasil.Err = ErrPermission
		cache.Policy = authorization_iface.Deny
	} else {
		for _, policy := range actualpermissions {
			if policy.Policy == authorization_iface.Deny {
				hasil.Err = ErrPermission
				cache.Policy = authorization_iface.Deny
				break
			}
		}
	}

	// saving cache
	err = auth.cache.Add(context.Background(), &ware_cache.CacheItem{
		Key:        cache.ID,
		Expiration: time.Hour * 6,
		Data:       cache,
	})
	// err = auth.tx.Transaction(func(tx *gorm.DB) error {
	// 	err = auth.tx.Save(&cache).Error
	// 	if err != nil {
	// 		return err
	// 	}

	// 	roles := actualpermissions.GetRole()
	// 	for _, role := range roles {
	// 		roleCache := RoleCache{
	// 			CacheAuthorizeID: cache.ID,
	// 			RoleID:           role.ID,
	// 		}
	// 		err = tx.Save(&roleCache).Error
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	return nil
	// })

	if err != nil {
		hasil.Err = ErrPermission
		return iscache, &hasil
	}

	if hasil.Err != nil {
		return iscache, &hasil
	}

	return iscache, nil

}

func NewAuthorization(cache ware_cache.Cache, tx *gorm.DB, passphrase string) *Authorization {
	// cache := ware_cache.NewMemcache()

	return &Authorization{
		tx:         tx,
		passphrase: passphrase,
		cache:      cache,
	}
}
