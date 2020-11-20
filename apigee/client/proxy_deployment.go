package client

import "strings"

const (
	ProxyDeploymentPath         = "o/%s/e/%s/apis/%s/deployments"
	ProxyDeploymentRevisionPath = "o/%s/e/%s/apis/%s/revisions/%d/deployments"
	ProxyDeploymentIdSeparator  = ":"
)

type ProxyDeployment struct {
	ProxyName       string                    `json:"name"`
	EnvironmentName string                    `json:"environment"`
	Revisions       []ProxyRevisionDeployment `json:"revision"`
}

type ProxyRevisionDeployment struct {
	Name string `json:"name"`
}

func (c *ProxyDeployment) ProxyDeploymentEncodeId() string {
	return c.EnvironmentName + ProxyDeploymentIdSeparator + c.ProxyName
}

func ProxyDeploymentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, ProxyDeploymentIdSeparator)
	return tokens[0], tokens[1]
}
