package authorization_iface

type Authorization interface {
	ApiQueryCheckPermission(identity Identity, query PermissionQuery) (bool, error)
}

type UserSystemCreate interface {
	Create() error
}

type AuthorizationMutation interface {
	UserSystemCreate(username string, apppwd string) UserSystemCreate
}

type SessionManager interface {
	FlushSession() error
	Session() Session
}

type Session interface {
	Create() Session
	Get() Session
	Delete() Session
	Err() error
}
