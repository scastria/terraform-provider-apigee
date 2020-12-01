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
