package identity

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pdcgo/shared/authorization"
	"github.com/pdcgo/shared/interfaces/identity_iface"
)

// type Agent interface {
// 	GetAgentType() string
// 	GetUserID() uint
// }

type ImporterAgent struct {
	*authorization.JwtIdentity
}

// GetAgentType implements Agent.
func (i *ImporterAgent) GetAgentType() identity_iface.AgentType {
	if i.JwtIdentity.UserAgent != "" {
		return i.JwtIdentity.UserAgent
	}

	return identity_iface.ImporterAgent
}

// GetUserID implements Agent.
// Subtle: this method shadows the method (*JwtIdentity).GetUserID of ImporterIdentity.JwtIdentity.
func (i *ImporterAgent) GetUserID() uint {
	return i.JwtIdentity.UserID
}

func NewEmptyImporterAgent() *ImporterAgent {
	return &ImporterAgent{
		JwtIdentity: &authorization.JwtIdentity{},
	}
}

type ApiAgent struct {
	*authorization.JwtIdentity
}

// GetToken implements identity_iface.Agent.
func (a *ApiAgent) GetToken(appname string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.JwtIdentity)
	tokenstring, err := token.SignedString([]byte(secret))
	return tokenstring, err
}

// GetAgentType implements Agent.
func (a *ApiAgent) GetAgentType() identity_iface.AgentType {
	if a.JwtIdentity.UserAgent != "" {
		return a.UserAgent
	}
	return "api"
}

// GetUserID implements Agent.
// Subtle: this method shadows the method (*JwtIdentity).GetUserID of ApiAgent.JwtIdentity.
func (a *ApiAgent) GetUserID() uint {
	return a.JwtIdentity.UserID
}

func NewEmptyApiAgent() *ApiAgent {
	return &ApiAgent{
		JwtIdentity: &authorization.JwtIdentity{},
	}
}
