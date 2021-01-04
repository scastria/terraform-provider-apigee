package client

import "strings"

const (
	ProxyDeploymentPath         = "organizations/%s/environments/%s/apis/%s/deployments"
	ProxyDeploymentRevisionPath = "organizations/%s/environments/%s/apis/%s/revisions/%d/deployments"
)

type ProxyDeployment struct {
	ProxyName       string                    `json:"name"`
	EnvironmentName string                    `json:"environment"`
	Revisions       []ProxyRevisionDeployment `json:"revision"`
}
type ProxyRevisionDeployment struct {
	Name string `json:"name"`
}
type GoogleProxyDeployment struct {
	Deployments []GoogleProxyDeploymentDeployments `json:"deployments"`
}
type GoogleProxyDeploymentDeployments struct {
	ProxyName       string `json:"apiProxy"`
	EnvironmentName string `json:"environment"`
	Revision        string `json:"revision"`
}

func (c *ProxyDeployment) ProxyDeploymentEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.ProxyName
}

func ProxyDeploymentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
