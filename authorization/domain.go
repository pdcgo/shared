package authorization

import (
	"fmt"
	"strings"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"gorm.io/gorm"
)

type Domain struct {
	ID   uint
	auth *Authorization
}

func (domain *Domain) CreateRole(name string) (*authorization_iface.Role, error) {
	var err error
	hasil := authorization_iface.Role{
		Key:      name,
		DomainID: domain.ID,
	}

	err = domain.auth.tx.Save(&hasil).Error
	if err != nil {
		err = domain.auth.tx.Where("key = ?", name).Where("domain_id = ?", domain.ID).First(&hasil).Error
	}

	return &hasil, err
}
func (domain *Domain) ListRole() ([]*authorization_iface.Role, error) {
	var err error
	hasil := []*authorization_iface.Role{}

	err = domain.auth.tx.Model(&authorization_iface.Role{}).
		Preload("Permissions").
		Where("domain_id = ?", domain.ID).
		Find(&hasil).
		Error

	return hasil, err
}

func (domain *Domain) DeleteRole(roleID uint) error {
	var err error
	err = domain.auth.refreshCacheRole(roleID)
	if err != nil {
		return err
	}

	err = domain.auth.tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&authorization_iface.Permission{}).
			Where("role_id = ?", roleID).
			Delete(&authorization_iface.Permission{}).
			Error

		if err != nil {
			return err
		}

		err = tx.Model(&authorization_iface.UserRole{}).
			Where("role_id = ?", roleID).
			Delete(&authorization_iface.UserRole{}).
			Error

		if err != nil {
			return err
		}

		err := tx.Model(&authorization_iface.Role{}).
			Where("id = ?", roleID).
			Delete(&authorization_iface.Role{}).
			Error

		return err
	})

	return err
}

//	type PermissionPayload struct {
//		EntityID string `json:"entity_id"`
//		Action   Action `json:"action"`
//		Policy   Policy `json:"policy"`
//	}
type AddPermItem struct {
	Action authorization_iface.Action `json:"action"`
	Policy authorization_iface.Policy `json:"policy"`
}

func (perm AddPermItem) ToTs(ident int) string {
	tabs := strings.Repeat("\t", ident)
	return fmt.Sprintf("%s%s?: boolean", tabs, perm.Action)
}

type AddPermissionPayload map[Entity][]*AddPermItem

func (perm AddPermissionPayload) ToTs(ident int) string {
	objectd := []string{}
	objectd = append(objectd, "{")

	tabs := strings.Repeat("\t", ident)

	for ent, perms := range perm {
		if len(perms) == 0 {
			continue
		}

		cident := ident + 1
		ctab := strings.Repeat("\t", cident)

		values := []string{}
		values = append(values, "{")
		for _, item := range perms {
			values = append(values, item.ToTs(cident+1))
		}
		values = append(values, tabs+"\t}")
		cvalue := strings.Join(values, "\n")

		valuedata := fmt.Sprintf("%s%s: %s", ctab, ent.GetEntityID(), cvalue)
		objectd = append(objectd, valuedata)
	}

	objectd = append(objectd, tabs+"}")
	hasil := strings.Join(objectd, "\n")
	return hasil
}

func (domain *Domain) RoleAddPermission(roleID uint, payload AddPermissionPayload) error {
	var err error
	err = domain.auth.refreshCacheRole(roleID)
	if err != nil {
		return err
	}

	err = domain.auth.tx.Transaction(func(tx *gorm.DB) error {
		for ent, item := range payload {
			for _, pol := range item {
				perm := authorization_iface.Permission{
					RoleID:   roleID,
					DomainID: domain.ID,
					EntityID: ent.GetEntityID(),
					Action:   pol.Action,
					Policy:   pol.Policy,
				}

				err := tx.Model(&authorization_iface.Permission{}).Where(&authorization_iface.Permission{
					RoleID:   roleID,
					DomainID: domain.ID,
					EntityID: ent.GetEntityID(),
					Action:   pol.Action,
					Policy:   pol.Policy,
				}).First(&perm).Error
				if err == nil {
					continue
				}

				err = tx.Save(&perm).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (domain *Domain) RoleListPermission(roleID uint) ([]*authorization_iface.Permission, error) {
	var err error
	hasil := []*authorization_iface.Permission{}

	err = domain.auth.tx.Model(&authorization_iface.Permission{}).
		Where("role_id = ?", roleID).
		Where("domain_id = ?", domain.ID).
		Find(&hasil).
		Error

	return hasil, err
}

func (domain *Domain) RoleRemovePermission(roleID uint, payload AddPermissionPayload) error {
	var err error

	for ent, item := range payload {
		for _, pol := range item {
			err = domain.auth.tx.
				Model(&authorization_iface.Permission{}).
				Where("role_id = ?", roleID).
				Where("domain_id = ?", domain.ID).
				Where("entity_id = ?", ent.GetEntityID()).
				Where("action = ?", pol.Action).
				Where("policy = ?", pol.Policy).
				Delete(&authorization_iface.Permission{}).
				Error

			if err != nil {
				return err
			}
		}
	}

	return err
}

func (domain *Domain) UserAddRole(identity authorization_iface.Identity, roleID uint) error {
	var err error
	userrole := authorization_iface.UserRole{
		RoleID: roleID,
		UserID: identity.IdentityID(),
	}
	err = domain.auth.tx.Save(&userrole).Error
	if err != nil {
		return err
	}
	err = domain.auth.RefreshCacheUser(identity.IdentityID())
	return err
}
func (domain *Domain) UserRemoveRole(userID uint, roleID uint) error {
	var err error

	err = domain.auth.tx.Model(&authorization_iface.UserRole{}).
		Where("role_id = ?", roleID).
		Where("user_id = ?", userID).
		Delete(&authorization_iface.UserRole{}).
		Error

	if err != nil {
		return err
	}

	err = domain.auth.RefreshCacheUser(userID)
	return err
}

func (domain *Domain) UserListRole(userID uint) ([]*authorization_iface.Role, error) {
	var err error
	hasil := []*authorization_iface.Role{}
	sub := domain.auth.tx.
		Model(&authorization_iface.UserRole{}).
		Select("role_id").
		Where("user_id = ?", userID)

	err = domain.auth.tx.Model(&authorization_iface.Role{}).
		Where("domain_id = ?", domain.ID).
		Where("id IN (?)", sub).
		Preload("Permissions").
		Find(&hasil).
		Error

	return hasil, err
}
