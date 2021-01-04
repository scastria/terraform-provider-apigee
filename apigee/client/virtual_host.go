package client

import "strings"

const (
	VirtualHostPath    = "organizations/%s/environments/%s/virtualhosts"
	VirtualHostPathGet = VirtualHostPath + "/%s"
)

type VirtualHost struct {
	EnvironmentName string
	Name            string   `json:"name"`
	HostAliases     []string `json:"hostAliases"`
	Port            string   `json:"port,omitempty"`
	BaseURL         string   `json:"baseUrl,omitempty"`
}

func (c *VirtualHost) VirtualHostEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func VirtualHostDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
