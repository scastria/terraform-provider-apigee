package client

import (
	"strings"
)

const (
	RolePermissionPath        = "o/%s/userroles/%s/permissions"
	RolePermissionPathGet     = RolePermissionPath + "/%s"
	RolePermissionIdSeparator = ":"
)

type RolePermission struct {
	RoleName    string
	Path        string   `json:"path"`
	Permissions []string `json:"permissions"`
}

func (rp *RolePermission) RolePermissionEncodeId() string {
	return rp.RoleName + RolePermissionIdSeparator + rp.Path
}

func RolePermissionDecodeId(s string) (string, string) {
	tokens := strings.Split(s, RolePermissionIdSeparator)
	return tokens[0], tokens[1]
}
