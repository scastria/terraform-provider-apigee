package client

const (
	ProxyPath    = "o/%s/apis"
	ProxyPathGet = ProxyPath + "/%s"
)

type ProxyRevision struct {
	Name       string `json:"name"`
	Revision   string `json:"revision"`
	Bundle     string
	BundleHash string
}

type Proxy struct {
	Name      string   `json:"name"`
	Revisions []string `json:"revision"`
}
