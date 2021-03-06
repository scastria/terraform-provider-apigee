package client

import "strings"

const (
	SharedFlowDeploymentPath         = "organizations/%s/environments/%s/sharedflows/%s/deployments"
	SharedFlowDeploymentRevisionPath = "organizations/%s/environments/%s/sharedflows/%s/revisions/%d/deployments"
)

type SharedFlowDeployment struct {
	SharedFlowName  string                         `json:"name"`
	EnvironmentName string                         `json:"environment"`
	Revisions       []SharedFlowRevisionDeployment `json:"revision"`
}
type SharedFlowRevisionDeployment struct {
	Name string `json:"name"`
}
type GoogleSharedFlowDeployment struct {
	Deployments []GoogleSharedFlowDeploymentDeployments `json:"deployments"`
}
type GoogleSharedFlowDeploymentDeployments struct {
	//Google API seems to reuse the structure from proxies for shared flows
	//Therefore, the use of apiProxy json property name is correct
	SharedFlowName  string `json:"apiProxy"`
	EnvironmentName string `json:"environment"`
	Revision        string `json:"revision"`
}

func (c *SharedFlowDeployment) SharedFlowDeploymentEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.SharedFlowName
}

func SharedFlowDeploymentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
