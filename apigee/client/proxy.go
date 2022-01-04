package client

const (
	ProxyPath           = "organizations/%s/apis"
	ProxyPathGet        = ProxyPath + "/%s"
	ProxyDeploymentPath = "organizations/%s/apis/%s/deployments"
)

type ProxyRevision struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
}

type Proxy struct {
	Name      string   `json:"name"`
	Revisions []string `json:"revision"`
}

type ProxyDeployments struct {
	Environments []ProxyDeployment `json:"environment"`
	ProxyName    string            `json:"name"`
}
type ProxyDeployment struct {
	EnvironmentName string `json:"name"`
}
