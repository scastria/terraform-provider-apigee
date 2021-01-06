package client

import "strings"

const (
	CompanyDeveloperPath    = "organizations/%s/companies/%s/developers"
	CompanyDeveloperPathGet = CompanyDeveloperPath + "/%s"
)

type CompanyDeveloper struct {
	CompanyName    string `json:"-"`
	DeveloperEmail string `json:"email"`
	Role           string `json:"role"`
}

type CompanyDeveloperList struct {
	Developers []CompanyDeveloper `json:"developer"`
}

func (ur *CompanyDeveloper) CompanyDeveloperEncodeId() string {
	return ur.CompanyName + IdSeparator + ur.DeveloperEmail
}

func CompanyDeveloperDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
