package client

import "strings"

const (
	TargetServerPath    = "organizations/%s/environments/%s/targetservers"
	TargetServerPathGet = TargetServerPath + "/%s"
)

type TargetServer struct {
	EnvironmentName string
	Name            string `json:"name"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	IsEnabled       bool   `json:"isEnabled"`
	SSLInfo         *SSL   `json:"sSLInfo,omitempty"`
}

type SSL struct {
	Enabled string `json:"enabled"`
}

func (c *TargetServer) TargetServerEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func TargetServerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
