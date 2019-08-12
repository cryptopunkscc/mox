package acl

type Permissions struct {
	Access       bool
	CheckBalance bool
}

type ACL struct {
	acl map[string]Permissions
}

func New() *ACL {
	acl := &ACL{
		acl: make(map[string]Permissions),
	}

	return acl
}

func (acl *ACL) Get(user string) Permissions {
	if p, ok := acl.acl[user]; ok {
		return p
	} else {
		return Permissions{}
	}
}

func (acl *ACL) Set(user string, p Permissions) {
	acl.acl[user] = p
}
