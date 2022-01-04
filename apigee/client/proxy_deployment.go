package client

import "strings"

const (
	ProxyEnvironmentDeploymentPath         = "organizations/%s/environments/%s/apis/%s/deployments"
	ProxyEnvironmentDeploymentRevisionPath = "organizations/%s/environments/%s/apis/%s/revisions/%d/deployments"
)

type ProxyEnvironmentDeployment struct {
	ProxyName       string                               `json:"name"`
	EnvironmentName string                               `json:"environment"`
	Revisions       []ProxyEnvironmentRevisionDeployment `json:"revision"`
}
type ProxyEnvironmentRevisionDeployment struct {
	Name string `json:"name"`
}
type GoogleProxyEnvironmentDeployment struct {
	Deployments []GoogleProxyEnvironmentDeploymentDeployments `json:"deployments"`
}
type GoogleProxyEnvironmentDeploymentDeployments struct {
	ProxyName       string `json:"apiProxy"`
	EnvironmentName string `json:"environment"`
	Revision        string `json:"revision"`
}

func (c *ProxyEnvironmentDeployment) ProxyDeploymentEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.ProxyName
}

func ProxyDeploymentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
