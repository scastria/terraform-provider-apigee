package client

import "strings"

const (
	DeveloperAppPath             = "o/%s/developers/%s/apps"
	DeveloperAppPathGet          = DeveloperAppPath + "/%s"
	DeveloperAppPathGeneratedKey = DeveloperAppPathGet + "/keys/%s"
	DeveloperAppIdSeparator      = ":"
)

type DeveloperApp struct {
	DeveloperEmail string
	Name           string                   `json:"name"`
	CallbackURL    string                   `json:"callbackUrl"`
	Attributes     []Attribute              `json:"attributes,omitempty"`
	Credentials    []DeveloperAppCredential `json:"credentials"`
}

func (ur *DeveloperApp) DeveloperAppEncodeId() string {
	return ur.DeveloperEmail + DeveloperAppIdSeparator + ur.Name
}

func DeveloperAppDecodeId(s string) (string, string) {
	tokens := strings.Split(s, DeveloperAppIdSeparator)
	return tokens[0], tokens[1]
}
