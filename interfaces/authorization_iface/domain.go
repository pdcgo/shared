package authorization_iface

type RoleAddPermissionItem struct {
	Action Action `json:"action"`
	Policy Policy `json:"policy"`
}

type RoleAddPermissionPayload map[Entity][]*RoleAddPermissionItem

type DomainV2 interface {
	RoleAddPermission(rolekey string, pay RoleAddPermissionPayload) error
	RoleAddPermissionWithDomain(rolekey string, domainID uint, pay RoleAddPermissionPayload) error
}
