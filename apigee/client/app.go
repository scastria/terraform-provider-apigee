package client

import "strings"

const (
	DeveloperAppPath             = "organizations/%s/developers/%s/apps"
	DeveloperAppPathGet          = DeveloperAppPath + "/%s"
	DeveloperAppPathGeneratedKey = DeveloperAppPathGet + "/keys/%s"
	CompanyAppPath               = "organizations/%s/companies/%s/apps"
	CompanyAppPathGet            = CompanyAppPath + "/%s"
	CompanyAppPathGeneratedKey   = CompanyAppPathGet + "/keys/%s"
)

type App struct {
	Name        string          `json:"name"`
	CallbackURL string          `json:"callbackUrl"`
	Attributes  []Attribute     `json:"attributes,omitempty"`
	Credentials []AppCredential `json:"credentials"`
	//Only used for developer context
	DeveloperEmail string `json:"-"`
	//Only used for company context
	CompanyName string `json:"-"`
}

func (ur *App) DeveloperAppEncodeId() string {
	return ur.DeveloperEmail + IdSeparator + ur.Name
}

func (ur *App) CompanyAppEncodeId() string {
	return ur.CompanyName + IdSeparator + ur.Name
}

func AppDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
