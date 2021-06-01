package client

import (
	"strings"
)

const (
	AliasPath    = "organizations/%s/environments/%s/keystores/%s/aliases"
	AliasPathGet = AliasPath + "/%s"
)

type Alias struct {
	EnvironmentName         string
	KeystoreName            string
	Name                    string
	Format                  string
	IgnoreExpiryValidation  bool
	IgnoreNewlineValidation bool
}

func (c *Alias) AliasEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.KeystoreName + IdSeparator + c.Name
}

func AliasDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
