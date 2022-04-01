package client

import "strings"

const (
	TargetServerPath    = "organizations/%s/environments/%s/targetservers"
	TargetServerPathGet = TargetServerPath + "/%s"
)

type TargetServer struct {
	EnvironmentName string `json:"-"`
	Name            string `json:"name"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	IsEnabled       bool   `json:"isEnabled"`
	SSLInfo         *SSL   `json:"sSLInfo,omitempty"`
}
type GoogleTargetServer struct {
	EnvironmentName string     `json:"-"`
	Name            string     `json:"name"`
	Host            string     `json:"host,omitempty"`
	Port            int        `json:"port,omitempty"`
	IsEnabled       bool       `json:"isEnabled"`
	SSLInfo         *GoogleSSL `json:"sSLInfo,omitempty"`
}

func (c *TargetServer) TargetServerEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}
func (c *GoogleTargetServer) TargetServerEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func TargetServerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
