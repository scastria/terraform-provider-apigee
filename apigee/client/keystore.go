package client

import "strings"

const (
	KeystorePath    = "organizations/%s/environments/%s/keystores"
	KeystorePathGet = KeystorePath + "/%s"
)

type Keystore struct {
	EnvironmentName string `json:"-"`
	Name            string `json:"name"`
}

func (c *Keystore) KeystoreEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func KeystoreDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
