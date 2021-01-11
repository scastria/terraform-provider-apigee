package client

type SSL struct {
	Enabled  string `json:"enabled"`
	KeyStore string `json:"keyStore,omitempty"`
	KeyAlias string `json:"keyAlias,omitempty"`
}
