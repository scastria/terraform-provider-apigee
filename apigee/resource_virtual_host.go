package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
	"mime"
	"net/http"
	"strconv"
)

func resourceVirtualHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualHostCreate,
		ReadContext:   resourceVirtualHostRead,
		UpdateContext: resourceVirtualHostUpdate,
		DeleteContext: resourceVirtualHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_aliases": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"base_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
			},
			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_keystore": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_keyalias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_truststore": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_client_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_ignore_validation_errors": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVirtualHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	aliasSet := d.Get("host_aliases").(*schema.Set)
	aliasList := convertSetToArray(aliasSet)
	if len(aliasList) == 0 {
		d.SetId("")
		return diag.Errorf("host_aliases must contain at least 1 value")
	}
	newVirtualHost := client.VirtualHost{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
		HostAliases:     aliasList,
	}
	fillVirtualHost(&newVirtualHost, d)
	err := json.NewEncoder(&buf).Encode(newVirtualHost)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.VirtualHostPath, c.Organization, newVirtualHost.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newVirtualHost.VirtualHostEncodeId())
	return diags
}

func fillVirtualHost(c *client.VirtualHost, d *schema.ResourceData) {
	port, ok := d.GetOk("port")
	if ok {
		c.Port = strconv.Itoa(port.(int))
	}
	baseURL, ok := d.GetOk("base_url")
	if ok {
		c.BaseURL = baseURL.(string)
	}
	sslEnabled, ok := d.GetOk("ssl_enabled")
	c.SSLInfo = &client.SSL{
		Enabled: strconv.FormatBool(sslEnabled.(bool)),
	}
	sslKeyStore, ok := d.GetOk("ssl_keystore")
	if ok {
		c.SSLInfo.KeyStore = sslKeyStore.(string)
	}
	sslKeyAlias, ok := d.GetOk("ssl_keyalias")
	if ok {
		c.SSLInfo.KeyAlias = sslKeyAlias.(string)
	}
	sslTrustStore, ok := d.GetOk("ssl_truststore")
	if ok {
		c.SSLInfo.TrustStore = sslTrustStore.(string)
	}
	sslClientAuthEnabled, ok := d.GetOk("ssl_client_auth_enabled")
	if ok {
		c.SSLInfo.ClientAuthEnabled = sslClientAuthEnabled.(bool)
	}
	sslIgnoreValidationErrors, ok := d.GetOk("ssl_ignore_validation_errors")
	if ok {
		c.SSLInfo.IgnoreValidationErrors = sslIgnoreValidationErrors.(bool)
	}
}

func resourceVirtualHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.VirtualHostDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.VirtualHostPathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.VirtualHost{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("host_aliases", retVal.HostAliases)
	port, _ := strconv.Atoi(retVal.Port)
	d.Set("port", port)
	d.Set("base_url", retVal.BaseURL)
	if retVal.SSLInfo != nil {
		sslEnabled, _ := strconv.ParseBool(retVal.SSLInfo.Enabled)
		d.Set("ssl_enabled", sslEnabled)
		d.Set("ssl_keystore", retVal.SSLInfo.KeyStore)
		d.Set("ssl_keyalias", retVal.SSLInfo.KeyAlias)
		d.Set("ssl_truststore", retVal.SSLInfo.TrustStore)
		d.Set("ssl_client_auth_enabled", retVal.SSLInfo.ClientAuthEnabled)
		d.Set("ssl_ignore_validation_errors", retVal.SSLInfo.IgnoreValidationErrors)
	} else {
		d.Set("ssl_enabled", false)
		d.Set("ssl_keystore", "")
		d.Set("ssl_keyalias", "")
		d.Set("ssl_truststore", "")
		d.Set("ssl_client_auth_enabled", false)
		d.Set("ssl_ignore_validation_errors", false)
	}
	return diags
}

func resourceVirtualHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.VirtualHostDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	aliasSet := d.Get("host_aliases").(*schema.Set)
	aliasList := convertSetToArray(aliasSet)
	if len(aliasList) == 0 {
		return diag.Errorf("host_aliases must contain at least 1 value")
	}
	upVirtualHost := client.VirtualHost{
		EnvironmentName: envName,
		Name:            name,
		HostAliases:     aliasList,
	}
	fillVirtualHost(&upVirtualHost, d)
	err := json.NewEncoder(&buf).Encode(upVirtualHost)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.VirtualHostPathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceVirtualHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.VirtualHostDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.VirtualHostPathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
