package client

import "strings"

const (
	AliasPath       = "organizations/%s/environments/%s/keystores/%s/aliases"
	AliasPathGet    = AliasPath + "/%s"
	AliasPathUpdate = AliasPath + "/%s"
)

func GetSupportedFormats() []string {
	return []string{"keycertfile", "keycertjar", "pkcs12"}
}

type Alias struct {
	EnvironmentName         string `json:"-"`
	Name                    string `json:"name"`
	CertFile                string `json:"cert_file"`
	File                    string `json:"file"`
	KeyFile                 string `json:"key_file"`
	KeystoreName            string `json:"keystore_name"`
	Format                  string `json:"format"`
	IgnoreExpiryValidation  bool   `json:"ignore_expiry_validation"`
	IgnoreNewlineValidation bool   `json:"ignore_newline_validation"`
}

func (a *Alias) AliasEncodeId() string {
	return a.KeystoreName + IdSeparator + a.Name
}

func AliasDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
