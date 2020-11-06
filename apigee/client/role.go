package client

const (
	RolePath    = "o/%s/userroles"
	RolePathGet = RolePath + "/%s"
)

type Role struct {
	Name string `json:"name"`
}

type RoleList struct {
	Roles []Role `json:"role"`
}
