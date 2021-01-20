package client

type SSL struct {
	Enabled                string `json:"enabled"`
	KeyStore               string `json:"keyStore,omitempty"`
	KeyAlias               string `json:"keyAlias,omitempty"`
	TrustStore             string `json:"trustStore,omitempty"`
	ClientAuthEnabled      string `json:"clientAuthEnabled,omitempty"`
	IgnoreValidationErrors bool   `json:"ignoreValidationErrors,omitempty"`
}
