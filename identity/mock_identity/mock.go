package mock_identity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pdcgo/shared/interfaces/identity_iface"
)

type MockJwt struct {
	jwt.StandardClaims
	UserID     uint
	ValidUntil int64
}

type MockAgent struct {
	UserID uint

	Type string
}

// GetToken implements identity_iface.Agent.
func (m *MockAgent) GetToken(appname string, secret string) (string, error) {
	validUntil := time.Now().Add(time.Hour * 350).UnixMicro()
	j := MockJwt{
		UserID:     m.UserID,
		ValidUntil: validUntil,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: validUntil,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	return token.SignedString([]byte(secret))
}

// GetAgentType implements identity_iface.Agent.
func (m *MockAgent) GetAgentType() identity_iface.AgentType {
	return identity_iface.AgentType(m.Type)
}

// GetUserID implements identity_iface.Agent.
func (m *MockAgent) GetUserID() uint {
	return m.UserID
}

func NewMockAgent(userid uint, tipe string) identity_iface.Agent {
	return &MockAgent{
		UserID: userid,
		Type:   tipe,
	}
}
