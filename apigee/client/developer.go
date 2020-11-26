package client

const (
	DeveloperPath    = "o/%s/developers"
	DeveloperPathGet = DeveloperPath + "/%s"
)

type Developer struct {
	Email      string               `json:"email"`
	FirstName  string               `json:"firstName"`
	LastName   string               `json:"lastName"`
	UserName   string               `json:"userName"`
	Attributes []DeveloperAttribute `json:"attributes,omitempty"`
}

type DeveloperAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
