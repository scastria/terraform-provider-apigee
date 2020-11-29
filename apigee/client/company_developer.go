package client

import "strings"

const (
	CompanyDeveloperPath        = "o/%s/companies/%s/developers"
	CompanyDeveloperPathGet     = CompanyDeveloperPath + "/%s"
	CompanyDeveloperIdSeparator = ":"
)

type CompanyDeveloper struct {
	CompanyName    string
	DeveloperEmail string `json:"email"`
	Role           string `json:"role"`
}

type CompanyDeveloperList struct {
	Developers []CompanyDeveloper `json:"developer"`
}

func (ur *CompanyDeveloper) CompanyDeveloperEncodeId() string {
	return ur.CompanyName + CompanyDeveloperIdSeparator + ur.DeveloperEmail
}

func CompanyDeveloperDecodeId(s string) (string, string) {
	tokens := strings.Split(s, CompanyDeveloperIdSeparator)
	return tokens[0], tokens[1]
}
