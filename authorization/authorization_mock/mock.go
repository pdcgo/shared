package authorization_mock

import (
	"testing"

	"github.com/pdcgo/gudang/src/authorization"
	"github.com/pdcgo/gudang/src/user"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/pkg/ware_cache"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type DomainScen struct {
	t    *testing.T
	D    *authorization.Domain
	auth *authorization.Authorization
}

type RoleScen struct {
	*DomainScen
	Role *authorization_iface.Role
}

func (d *RoleScen) WithUser(identity authorization_iface.Identity, handler func()) {
	err := d.D.UserAddRole(identity, d.Role.ID)
	assert.Nil(d.t, err)

	defer func() {
		err := d.D.UserRemoveRole(identity.IdentityID(), d.Role.ID)
		assert.Nil(d.t, err)
	}()

	handler()
}

type RolesScen struct {
	*DomainScen
	Roles map[string]*authorization_iface.Role
}

func (r *RolesScen) WithUser(key string, identity authorization_iface.Identity) {
	err := r.DomainScen.D.UserAddRole(identity, r.Roles[key].ID)
	assert.Nil(r.t, err)
}

func (d *DomainScen) SeedDefault(db *gorm.DB, tipe db_models.TeamType) *RolesScen {
	err := user.CreateDefaultRoleTeam(d.auth, d.D.ID, tipe)
	assert.Nil(d.t, err)
	hasil := map[string]*authorization_iface.Role{}

	roles, err := d.D.ListRole()
	assert.Nil(d.t, err)

	for _, role := range roles {
		d := role
		hasil[role.Key] = d
	}

	return &RolesScen{
		DomainScen: d,
		Roles:      hasil,
	}
}

func (d *DomainScen) WithRole(roleName string, perms authorization.AddPermissionPayload, handle func(role *RoleScen)) {
	role, err := d.D.CreateRole(roleName)
	assert.Nil(d.t, err)
	err = d.D.RoleAddPermission(role.ID, perms)
	assert.Nil(d.t, err)

	defer func() {
		err := d.D.RoleRemovePermission(role.ID, perms)
		assert.Nil(d.t, err)
		err = d.D.DeleteRole(role.ID)
		assert.Nil(d.t, err)

	}()
	assert.NotEqual(d.t, 0, role.ID)
	handle(&RoleScen{
		DomainScen: d,
		Role:       role,
	})

}

type AuthScen struct {
	t *testing.T
	A *authorization.Authorization
}

func (a *AuthScen) Domain(domainID uint) *DomainScen {
	return &DomainScen{
		t:    a.t,
		D:    a.A.Domain(domainID),
		auth: a.A,
	}
}

func NewAuthMock(t *testing.T, db *gorm.DB, handle func(auth *AuthScen)) {
	auth := authorization.NewAuthorization(ware_cache.NewLocalCache(), db, "test_phrase")
	handle(&AuthScen{
		t: t,
		A: auth,
	})
}
