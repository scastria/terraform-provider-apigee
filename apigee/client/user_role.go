package client

import "strings"

const (
	UserRolePath    = "o/%s/userroles/%s/users"
	UserRolePathGet = UserRolePath + "/%s"
)

type UserRole struct {
	EmailId  string
	RoleName string
}

func (ur *UserRole) UserRoleEncodeId() string {
	return ur.EmailId + IdSeparator + ur.RoleName
}

func UserRoleDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
