package client

import (
	"strconv"
	"strings"
)

const (
	ProxyPolicyPath    = "organizations/%s/apis/%s/revisions/%d/policies"
	ProxyPolicyPathGet = ProxyPolicyPath + "/%s"
)

type ProxyPolicy struct {
	ProxyName string
	Revision  int
	Name      string
}

func (c *ProxyPolicy) ProxyPolicyEncodeId() string {
	return c.ProxyName + IdSeparator + strconv.Itoa(c.Revision) + IdSeparator + c.Name
}

func ProxyPolicyDecodeId(s string) (string, int, string) {
	tokens := strings.Split(s, IdSeparator)
	revision, _ := strconv.Atoi(tokens[1])
	return tokens[0], revision, tokens[2]
}
