package client

import "strings"

const (
	AliasPath    = "organizations/%s/environments/%s/keystores/%s/aliases"
	AliasPathGet = AliasPath + "/%s"
)

func GetSupportedFormats() []string {
	return []string{"keycertfile", "keycertjar", "pkcs12"}
}

type Alias struct {
	EnvironmentName         string `json:"-"`
	Alias                   string `json:"alias"`
	CertFile                string `json:"cert_file"`
	File                    string `json:"file"`
	KeyFile                 string `json:"key_file"`
	KeystoreName            string `json:"keystore_name"`
	Format                  string `json:"format"`
	IgnoreExpiryValidation  bool   `json:"ignore_expiry_validation"`
	IgnoreNewlineValidation bool   `json:"ignore_newline_validation"`
}

func (a *Alias) AliasEncodeId() string {
	return a.EnvironmentName + IdSeparator + a.KeystoreName + IdSeparator + a.Alias
}

func AliasDecodeId(s string) (string, string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1], tokens[2]
}
