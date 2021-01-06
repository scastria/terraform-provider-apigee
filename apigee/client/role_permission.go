package client

import (
	"strings"
)

const (
	RolePermissionPath    = "organizations/%s/userroles/%s/permissions"
	RolePermissionPathGet = RolePermissionPath + "/%s"
)

type RolePermission struct {
	RoleName    string   `json:"-"`
	Path        string   `json:"path"`
	Permissions []string `json:"permissions"`
}

func (rp *RolePermission) RolePermissionEncodeId() string {
	return rp.RoleName + IdSeparator + rp.Path
}

func RolePermissionDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
