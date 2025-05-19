package identity

import "github.com/pdcgo/shared/interfaces/identity_iface"

type SystemAgent struct{}

// GetAgentType implements identity_iface.Agent.
func (s *SystemAgent) GetAgentType() identity_iface.AgentType {
	return identity_iface.SystemAgent
}

// GetUserID implements identity_iface.Agent.
func (s *SystemAgent) GetUserID() uint {
	return 1
}

func NewSystemAgent() *SystemAgent {
	return &SystemAgent{}
}
