package client

const (
	ProductPath        = "o/%s/apiproducts"
	ProductPathGet     = ProductPath + "/%s"
	AutoApprovalType   = "auto"
	ManualApprovalType = "manual"
)

type Product struct {
	APIResources  []string           `json:"apiResources,omitempty"`
	ApprovalType  string             `json:"approvalType,omitempty"`
	Attributes    []ProductAttribute `json:"attributes,omitempty"`
	Description   string             `json:"description,omitempty"`
	DisplayName   string             `json:"displayName,omitempty"`
	Environments  []string           `json:"environments,omitempty"`
	Name          string             `json:"name,omitempty"`
	Proxies       []string           `json:"proxies,omitempty"`
	Quota         string             `json:"quota,omitempty"`
	QuotaInterval string             `json:"quotaInterval,omitempty"`
	QuotaTimeUnit string             `json:"quotaTimeUnit,omitempty"`
	Scopes        []string           `json:"scopes,omitempty"`
}

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
