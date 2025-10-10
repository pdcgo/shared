package authorization

import (
	"errors"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/pkg/debugtool"
	"gorm.io/gorm"
)

type domainImpl struct {
	domainID uint
	db       *gorm.DB
}

// RoleAddPermission implements authorization_iface.DomainV2.
func (d *domainImpl) RoleAddPermission(rolekey string, payload authorization_iface.RoleAddPermissionPayload) error {
	var err error
	var role authorization_iface.Role
	err = d.
		db.
		Model(&authorization_iface.Role{}).
		Where("domain_id = ?", d.domainID).
		Where("key = ?", rolekey).
		First(&role).
		Error

	if err != nil {
		return errors.New("role not found")
	}

	for ent, item := range payload {
		for _, pol := range item {

			perm := authorization_iface.Permission{
				RoleID:   role.ID,
				DomainID: d.domainID,
				EntityID: ent.GetEntityID(),
				Action:   pol.Action,
				Policy:   pol.Policy,
			}

			err := d.db.
				Model(&authorization_iface.Permission{}).
				Where(&authorization_iface.Permission{
					RoleID:   role.ID,
					DomainID: d.domainID,
					EntityID: ent.GetEntityID(),
					Action:   pol.Action,
					Policy:   pol.Policy,
				}).
				First(&perm).
				Error

			if err == nil {
				// debugtool.LogJson(perm)
				continue
			}

			debugtool.LogJson(perm)

			err = d.db.Save(&perm).Error
			if err != nil {
				return err
			}
		}
	}

	return err
}

func NewDomainV2(db *gorm.DB, domainID uint) authorization_iface.DomainV2 {
	return &domainImpl{
		db:       db,
		domainID: domainID,
	}
}
