package client

const (
	DeveloperPath    = "organizations/%s/developers"
	DeveloperPathGet = DeveloperPath + "/%s"
)

type Developer struct {
	Email      string      `json:"email"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	UserName   string      `json:"userName"`
	Attributes []Attribute `json:"attributes,omitempty"`
}
