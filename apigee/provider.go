package apigee

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_PASSWORD", nil),
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_SERVER", nil),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("APIGEE_PORT", 443),
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_ORGANIZATION", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"apigee_user":                     resourceUser(),
			"apigee_role":                     resourceRole(),
			"apigee_user_role":                resourceUserRole(),
			"apigee_role_permission":          resourceRolePermission(),
			"apigee_cache":                    resourceCache(),
			"apigee_organization_kvm":         resourceOrganizationKVM(),
			"apigee_environment_kvm":          resourceEnvironmentKVM(),
			"apigee_proxy_kvm":                resourceProxyKVM(),
			"apigee_target_server":            resourceTargetServer(),
			"apigee_virtual_host":             resourceVirtualHost(),
			"apigee_proxy":                    resourceProxy(),
			"apigee_proxy_deployment":         resourceProxyDeployment(),
			"apigee_shared_flow":              resourceSharedFlow(),
			"apigee_shared_flow_deployment":   resourceSharedFlowDeployment(),
			"apigee_developer":                resourceDeveloper(),
			"apigee_product":                  resourceProduct(),
			"apigee_company":                  resourceCompany(),
			"apigee_company_developer":        resourceCompanyDeveloper(),
			"apigee_developer_app":            resourceDeveloperApp(),
			"apigee_developer_app_credential": resourceDeveloperAppCredential(),
			"apigee_company_app":              resourceCompanyApp(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"apigee_user": dataSourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	server := d.Get("server").(string)
	port := d.Get("port").(int)
	organization := d.Get("organization").(string)

	var diags diag.Diagnostics
	return client.NewClient(username, password, server, port, organization), diags
}
