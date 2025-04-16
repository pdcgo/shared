package identity_iface

type AgentType string

const (
	ApiAgent      AgentType = "api"
	ImporterAgent AgentType = "importer_api"
	ThirdAppAgent AgentType = "third_app"
	TestAgent     AgentType = "test_agent"
	SystemAgent   AgentType = "system"
)

type Agent interface {
	GetAgentType() AgentType
	GetUserID() uint
}
