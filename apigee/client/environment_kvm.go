package client

import "strings"

const (
	EnvironmentKVMPath         = "o/%s/e/%s/keyvaluemaps"
	EnvironmentKVMPathGet      = EnvironmentKVMPath + "/%s"
	EnvironmentKVMPathGetEntry = EnvironmentKVMPathGet + "/entries/%s"
	EnvironmentKVMIdSeparator  = ":"
)

type EnvironmentKVM struct {
	EnvironmentName string
	Name            string     `json:"name"`
	Encrypted       bool       `json:"encrypted,omitempty"`
	Entries         []KVMEntry `json:"entry,omitempty"`
}

type KVMEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *EnvironmentKVM) EnvironmentKVMEncodeId() string {
	return c.EnvironmentName + EnvironmentKVMIdSeparator + c.Name
}

func EnvironmentKVMDecodeId(s string) (string, string) {
	tokens := strings.Split(s, EnvironmentKVMIdSeparator)
	return tokens[0], tokens[1]
}
