package client

type SSLCommonName struct {
	Value         string `json:"value,omitempty"`
	WildcardMatch bool   `json:"wildcardMatch,omitempty"`
}

type SSL struct {
	Enabled                string         `json:"enabled"`
	KeyStore               string         `json:"keyStore,omitempty"`
	KeyAlias               string         `json:"keyAlias,omitempty"`
	TrustStore             string         `json:"trustStore,omitempty"`
	CommonName             *SSLCommonName `json:"commonName,omitempty"`
	ClientAuthEnabled      string         `json:"clientAuthEnabled,omitempty"`
	IgnoreValidationErrors bool           `json:"ignoreValidationErrors,omitempty"`
	Protocols              []string       `json:"protocols,omitempty"`
}
type GoogleSSL struct {
	Enabled                bool           `json:"enabled"`
	KeyStore               string         `json:"keyStore,omitempty"`
	KeyAlias               string         `json:"keyAlias,omitempty"`
	TrustStore             string         `json:"trustStore,omitempty"`
	CommonName             *SSLCommonName `json:"commonName,omitempty"`
	ClientAuthEnabled      bool           `json:"clientAuthEnabled,omitempty"`
	IgnoreValidationErrors bool           `json:"ignoreValidationErrors,omitempty"`
	Protocols              []string       `json:"protocols,omitempty"`
}
