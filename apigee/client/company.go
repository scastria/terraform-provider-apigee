package client

const (
	CompanyPath    = "o/%s/companies"
	CompanyPathGet = CompanyPath + "/%s"
)

type Company struct {
	Name        string             `json:"name"`
	DisplayName string             `json:"displayName"`
	Attributes  []CompanyAttribute `json:"attributes,omitempty"`
}

type CompanyAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
