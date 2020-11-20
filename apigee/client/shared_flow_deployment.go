package client

import "strings"

const (
	SharedFlowDeploymentPath         = "o/%s/e/%s/sharedflows/%s/deployments"
	SharedFlowDeploymentRevisionPath = "o/%s/e/%s/sharedflows/%s/revisions/%d/deployments"
	SharedFlowDeploymentIdSeparator  = ":"
)

type SharedFlowDeployment struct {
	SharedFlowName  string                         `json:"name"`
	EnvironmentName string                         `json:"environment"`
	Revisions       []SharedFlowRevisionDeployment `json:"revision"`
}

type SharedFlowRevisionDeployment struct {
	Name string `json:"name"`
}

func (c *SharedFlowDeployment) SharedFlowDeploymentEncodeId() string {
	return c.EnvironmentName + SharedFlowDeploymentIdSeparator + c.SharedFlowName
}

func SharedFlowDeploymentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, SharedFlowDeploymentIdSeparator)
	return tokens[0], tokens[1]
}
