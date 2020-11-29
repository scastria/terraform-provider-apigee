package client

import "strings"

const (
	OrganizationKVMPath         = "o/%s/keyvaluemaps"
	OrganizationKVMPathGet      = OrganizationKVMPath + "/%s"
	OrganizationKVMPathGetEntry = OrganizationKVMPathGet + "/entries/%s"
	EnvironmentKVMPath          = "o/%s/e/%s/keyvaluemaps"
	EnvironmentKVMPathGet       = EnvironmentKVMPath + "/%s"
	EnvironmentKVMPathGetEntry  = EnvironmentKVMPathGet + "/entries/%s"
	ProxyKVMPath                = "o/%s/apis/%s/keyvaluemaps"
	ProxyKVMPathGet             = ProxyKVMPath + "/%s"
	ProxyKVMPathGetEntry        = ProxyKVMPathGet + "/entries/%s"
	KVMIdSeparator              = ":"
)

type KVM struct {
	Name      string     `json:"name"`
	Encrypted bool       `json:"encrypted,omitempty"`
	Entries   []KVMEntry `json:"entry,omitempty"`
	//Only used for Environment context
	EnvironmentName string
	//Only used for Proxy context
	ProxyName string
}

type KVMEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *KVM) EnvironmentKVMEncodeId() string {
	return c.EnvironmentName + KVMIdSeparator + c.Name
}

func EnvironmentKVMDecodeId(s string) (string, string) {
	tokens := strings.Split(s, KVMIdSeparator)
	return tokens[0], tokens[1]
}

func (c *KVM) ProxyKVMEncodeId() string {
	return c.ProxyName + KVMIdSeparator + c.Name
}

func ProxyKVMDecodeId(s string) (string, string) {
	tokens := strings.Split(s, KVMIdSeparator)
	return tokens[0], tokens[1]
}
