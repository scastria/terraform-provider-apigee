package client

const (
	CompanyPath    = "organizations/%s/companies"
	CompanyPathGet = CompanyPath + "/%s"
)

type Company struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"displayName"`
	Attributes  []Attribute `json:"attributes,omitempty"`
}
