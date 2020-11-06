package client

type Role struct {
	Name string `json:"name"`
}

type RoleList struct {
	Roles []Role `json:"role"`
}
