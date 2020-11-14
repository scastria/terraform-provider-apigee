package client

import "strings"

const (
	TargetServerPath        = "o/%s/e/%s/targetservers"
	TargetServerPathGet     = TargetServerPath + "/%s"
	TargetServerIdSeparator = ":"
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
	return c.EnvironmentName + TargetServerIdSeparator + c.Name
}

func TargetServerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, TargetServerIdSeparator)
	return tokens[0], tokens[1]
}
