package client

const (
	SharedFlowPath    = "o/%s/sharedflows"
	SharedFlowPathGet = SharedFlowPath + "/%s"
)

type SharedFlowRevision struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
}

type SharedFlow struct {
	Name      string   `json:"name"`
	Revisions []string `json:"revision"`
}
