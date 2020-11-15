package client

import "strings"

const (
	VirtualHostPath        = "o/%s/e/%s/virtualhosts"
	VirtualHostPathGet     = VirtualHostPath + "/%s"
	VirtualHostIdSeparator = ":"
)

type VirtualHost struct {
	EnvironmentName string
	Name            string   `json:"name"`
	HostAliases     []string `json:"hostAliases"`
	Port            string   `json:"port,omitempty"`
	BaseURL         string   `json:"baseUrl,omitempty"`
}

func (c *VirtualHost) VirtualHostEncodeId() string {
	return c.EnvironmentName + VirtualHostIdSeparator + c.Name
}

func VirtualHostDecodeId(s string) (string, string) {
	tokens := strings.Split(s, VirtualHostIdSeparator)
	return tokens[0], tokens[1]
}
