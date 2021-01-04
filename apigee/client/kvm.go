package client

import "strings"

const (
	OrganizationKVMPath           = "organizations/%s/keyvaluemaps"
	OrganizationKVMPathGet        = OrganizationKVMPath + "/%s"
	OrganizationKVMPathEntries    = OrganizationKVMPathGet + "/entries"
	OrganizationKVMPathEntriesGet = OrganizationKVMPathEntries + "/%s"
	EnvironmentKVMPath            = "organizations/%s/environments/%s/keyvaluemaps"
	EnvironmentKVMPathGet         = EnvironmentKVMPath + "/%s"
	EnvironmentKVMPathEntries     = EnvironmentKVMPathGet + "/entries"
	EnvironmentKVMPathEntriesGet  = EnvironmentKVMPathEntries + "/%s"
	ProxyKVMPath                  = "organizations/%s/apis/%s/keyvaluemaps"
	ProxyKVMPathGet               = ProxyKVMPath + "/%s"
	ProxyKVMPathEntries           = ProxyKVMPathGet + "/entries"
	ProxyKVMPathEntriesGet        = ProxyKVMPathEntries + "/%s"
)

type KVM struct {
	Name      string      `json:"name"`
	Encrypted bool        `json:"encrypted,omitempty"`
	Entries   []Attribute `json:"entry,omitempty"`
	//Only used for Environment context
	EnvironmentName string
	//Only used for Proxy context
	ProxyName string
}

func (c *KVM) EnvironmentKVMEncodeId() string {
	return c.EnvironmentName + IdSeparator + c.Name
}

func (c *KVM) ProxyKVMEncodeId() string {
	return c.ProxyName + IdSeparator + c.Name
}

func KVMDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
