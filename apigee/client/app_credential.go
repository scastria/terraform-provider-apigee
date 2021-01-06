package client

import "strings"

const (
	DeveloperAppCredentialPath        = "organizations/%s/developers/%s/apps/%s/keys"
	DeveloperAppCredentialPathCreate  = DeveloperAppCredentialPath + "/create"
	DeveloperAppCredentialPathGet     = DeveloperAppCredentialPath + "/%s"
	DeveloperAppCredentialPathProduct = DeveloperAppCredentialPathGet + "/apiproducts/%s"
	CompanyAppCredentialPath          = "organizations/%s/companies/%s/apps/%s/keys"
	CompanyAppCredentialPathCreate    = CompanyAppCredentialPath + "/create"
	CompanyAppCredentialPathGet       = CompanyAppCredentialPath + "/%s"
	CompanyAppCredentialPathProduct   = CompanyAppCredentialPathGet + "/apiproducts/%s"
)

type AppCredential struct {
	AppName        string             `json:"-"`
	ConsumerKey    string             `json:"consumerKey"`
	ConsumerSecret string             `json:"consumerSecret"`
	Scopes         []string           `json:"scopes"`
	APIProducts    []APIProductStatus `json:"apiProducts"`
	Attributes     []Attribute        `json:"attributes,omitempty"`
	//Only used for developer context
	DeveloperEmail string `json:"-"`
	//Only used for company context
	CompanyName string `json:"-"`
}

type AppCredentialModify struct {
	AppName        string      `json:"-"`
	ConsumerKey    string      `json:"consumerKey"`
	ConsumerSecret string      `json:"consumerSecret"`
	Scopes         []string    `json:"scopes"`
	APIProducts    []string    `json:"apiProducts"`
	Attributes     []Attribute `json:"attributes,omitempty"`
	//Only used for developer context
	DeveloperEmail string `json:"-"`
	//Only used for company context
	CompanyName string `json:"-"`
}

type APIProductStatus struct {
	APIProduct string `json:"apiproduct"`
	Status     string `json:"status"`
}

func (ur *AppCredentialModify) DeveloperAppCredentialEncodeId() string {
	return ur.DeveloperEmail + IdSeparator + ur.AppName + IdSeparator + ur.ConsumerKey
}

func (ur *AppCredentialModify) CompanyAppCredentialEncodeId() string {
	return ur.CompanyName + IdSeparator + ur.AppName + IdSeparator + ur.ConsumerKey
}

func AppCredentialDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
