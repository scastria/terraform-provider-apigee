package client

import "strings"

const (
	UserRolePath        = "o/%s/userroles/%s/users"
	UserRolePathGet     = UserRolePath + "/%s"
	UserRoleIdSeparator = ":"
)

type UserRole struct {
	EmailId  string
	RoleName string
}

func (ur *UserRole) UserRoleEncodeId() string {
	return ur.EmailId + UserRoleIdSeparator + ur.RoleName
}

func UserRoleDecodeId(s string) (string, string) {
	tokens := strings.Split(s, UserRoleIdSeparator)
	return tokens[0], tokens[1]
}
