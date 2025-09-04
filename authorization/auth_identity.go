package authorization

import (
	"errors"
	"net/http"
	"strings"

	"github.com/pdcgo/shared/interfaces/authorization_iface"
)

type authIdentityImpl struct {
	identity authorization_iface.Identity
	auth     authorization_iface.Authorization
	err      error
}

// Err implements authorization_iface.AuthIdentity.
func (a *authIdentityImpl) Err() error {
	return a.err
}

// HasPermission implements authorization_iface.AuthIdentity.
func (a *authIdentityImpl) HasPermission(perms authorization_iface.CheckPermissionGroup) authorization_iface.AuthIdentity {
	var err error
	if a.err != nil {
		return a
	}

	// debugtool.LogJson(a.identity)

	if a.identity.IsSuperUser() {
		return a
	}

	err = a.auth.HasPermission(a.identity, perms)
	return a.setErr(err)
}

// Identity implements authorization_iface.AuthIdentity.
func (a *authIdentityImpl) Identity() authorization_iface.Identity {
	return a.identity
}

func (a *authIdentityImpl) parse(header http.Header, passphrase string) *authIdentityImpl {
	var err error
	token := header.Get("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")

	identity := JwtIdentity{}
	err = identity.Deserialize(passphrase, token)
	if err != nil {
		return a.setErr(errors.New("token expired or not authorization empty"))
	}

	a.identity = &identity
	return a
}

func (a *authIdentityImpl) setErr(err error) *authIdentityImpl {
	if a.err != nil {
		return a
	}

	if err != nil {
		a.err = err
	}

	return a
}

func NewAuthIdentityHttpHeader(auth authorization_iface.Authorization, header http.Header, passphrase string) authorization_iface.AuthIdentity {
	a := &authIdentityImpl{
		identity: &JwtIdentity{},
		auth:     auth,
	}

	a.parse(header, passphrase)

	return a
}
