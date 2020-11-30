package client

import "strings"

const (
	DeveloperAppCredentialPath        = "o/%s/developers/%s/apps/%s/keys"
	DeveloperAppCredentialPathCreate  = DeveloperAppCredentialPath + "/create"
	DeveloperAppCredentialPathGet     = DeveloperAppCredentialPath + "/%s"
	DeveloperAppCredentialPathProduct = DeveloperAppCredentialPathGet + "/apiproducts/%s"
)

type DeveloperAppCredential struct {
	DeveloperEmail   string
	DeveloperAppName string
	ConsumerKey      string             `json:"consumerKey"`
	ConsumerSecret   string             `json:"consumerSecret"`
	Scopes           []string           `json:"scopes"`
	APIProducts      []APIProductStatus `json:"apiProducts"`
	Attributes       []Attribute        `json:"attributes,omitempty"`
}

type DeveloperAppCredentialModify struct {
	DeveloperEmail   string
	DeveloperAppName string
	ConsumerKey      string      `json:"consumerKey"`
	ConsumerSecret   string      `json:"consumerSecret"`
	Scopes           []string    `json:"scopes"`
	APIProducts      []string    `json:"apiProducts"`
	Attributes       []Attribute `json:"attributes,omitempty"`
}

type APIProductStatus struct {
	APIProduct string `json:"apiproduct"`
	Status     string `json:"status"`
}

func (ur *DeveloperAppCredentialModify) DeveloperAppCredentialEncodeId() string {
	return ur.DeveloperEmail + IdSeparator + ur.DeveloperAppName + IdSeparator + ur.ConsumerKey
}

func DeveloperAppCredentialDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
