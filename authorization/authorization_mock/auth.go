package authorization_mock

import (
	"net/http"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/interfaces/identity_iface"
	"gorm.io/gorm"
)

type EmptyAuthorizationMock struct {
	AuthIdentityMock *AuthIdentityMock
}

// ApiQueryCheckPermission implements authorization_iface.Authorization.
func (e *EmptyAuthorizationMock) ApiQueryCheckPermission(
	identity authorization_iface.Identity,
	query authorization_iface.PermissionQuery,
) (bool, error) {
	return true, nil
}

// AuthIdentityFromHeader implements authorization_iface.Authorization.
func (e *EmptyAuthorizationMock) AuthIdentityFromHeader(header http.Header) authorization_iface.AuthIdentity {
	return e.AuthIdentityMock
}

// AuthIdentityFromToken implements authorization_iface.Authorization.
func (e *EmptyAuthorizationMock) AuthIdentityFromToken(token string) authorization_iface.AuthIdentity {
	return e.AuthIdentityMock
}

// HasPermission implements authorization_iface.Authorization.
func (e *EmptyAuthorizationMock) HasPermission(
	identity authorization_iface.Identity,
	perms authorization_iface.CheckPermissionGroup,
) error {
	return nil
}

type AuthIdentityMock struct {
	IdentityMock *IdentityMock
}

// Err implements authorization_iface.AuthIdentity.
func (m *AuthIdentityMock) Err() error {
	return nil
}

// HasPermission implements authorization_iface.AuthIdentity.
func (m *AuthIdentityMock) HasPermission(perms authorization_iface.CheckPermissionGroup) authorization_iface.AuthIdentity {
	return m
}

// Identity implements authorization_iface.AuthIdentity.
func (m *AuthIdentityMock) Identity() authorization_iface.Identity {
	return m.IdentityMock
}

type IdentityMock struct {
	ID uint
}

// GetAgentType implements authorization_iface.Identity.
func (i *IdentityMock) GetAgentType() identity_iface.AgentType {
	return "test"
}

// GetExpired implements authorization_iface.Identity.
func (i *IdentityMock) GetExpired(tx *gorm.DB) (*authorization_iface.ExpiredToken, error) {
	panic("unimplemented")
}

// GetToken implements authorization_iface.Identity.
func (i *IdentityMock) GetToken(appname string, secret string) (string, error) {
	return "", nil
}

// GetUserID implements authorization_iface.Identity.
func (i *IdentityMock) GetUserID() uint {
	return i.ID
}

// HasRole implements authorization_iface.Identity.
func (i *IdentityMock) HasRole(tx *gorm.DB, domainID uint, keyname string) (bool, error) {
	return true, nil
}

// IdentityID implements authorization_iface.Identity.
func (i *IdentityMock) IdentityID() uint {
	return i.ID
}

// IsSuperUser implements authorization_iface.Identity.
func (i *IdentityMock) IsSuperUser() bool {
	panic("unimplemented")
}

// IsTokenExpired implements authorization_iface.Identity.
func (i *IdentityMock) IsTokenExpired(tx *gorm.DB) (bool, error) {
	return false, nil
}
