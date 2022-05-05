package client

const (
	ProductPath        = "organizations/%s/apiproducts"
	ProductPathGet     = ProductPath + "/%s"
	AutoApprovalType   = "auto"
	ManualApprovalType = "manual"
)

type Product struct {
	APIResources   []string       `json:"apiResources,omitempty"`
	ApprovalType   string         `json:"approvalType,omitempty"`
	Attributes     []Attribute    `json:"attributes,omitempty"`
	Description    string         `json:"description,omitempty"`
	DisplayName    string         `json:"displayName,omitempty"`
	Environments   []string       `json:"environments,omitempty"`
	Name           string         `json:"name,omitempty"`
	OperationGroup OperationGroup `json:"operationGroup"`
	Proxies        []string       `json:"proxies,omitempty"`
	Quota          string         `json:"quota,omitempty"`
	QuotaInterval  string         `json:"quotaInterval,omitempty"`
	QuotaTimeUnit  string         `json:"quotaTimeUnit,omitempty"`
	Scopes         []string       `json:"scopes,omitempty"`
}

type OperationGroup struct {
	OperationConfigs    []OperationConfigs `json:"operationConfigs"`
	OperationConfigType string             `json:"operationConfigType"`
}

type Quota struct {
	Limit    string `json:"limit,omitempty"`
	Interval string `json:"interval,omitempty"`
	TimeUnit string `json:"timeUnit,omitempty"`
}

type Operation struct {
	Resource string   `json:"resource"`
	Methods  []string `json:"methods,omitempty"`
}

type OperationConfigs struct {
	ApiSource  string      `json:"apiSource"`
	Operations []Operation `json:"operations"`
	Quota      Quota       `json:"quota,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
}
