package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-apigee/apigee"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: apigee.Provider,
	})
}
