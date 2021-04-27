package client

import "strings"

const (
	ReferencePath    = "organizations/%s/environments/%s/references"
	ReferencePathGet = ReferencePath + "/%s"
)

type Reference struct {
	EnvironmentName string     `json:"-"`
	Name            string     `json:"name"`
	Refers            string     `json:"refers"`
	ResourceType            string     `json:"resourceType"`
	//OverflowToDisk                    bool       `json:"overflowToDisk,omitempty"`
}



func (c *Reference) ReferenceEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func ReferenceDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
