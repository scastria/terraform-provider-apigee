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

func resourceTargetServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTargetServerCreate,
		ReadContext:   resourceTargetServerRead,
		UpdateContext: resourceTargetServerUpdate,
		DeleteContext: resourceTargetServerDelete,
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
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

func resourceTargetServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newTargetServer := client.TargetServer{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
		Host:            d.Get("host").(string),
		Port:            d.Get("port").(int),
	}
	fillTargetServer(&newTargetServer, d)
	err := json.NewEncoder(&buf).Encode(newTargetServer)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.TargetServerPath, c.Organization, newTargetServer.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newTargetServer.TargetServerEncodeId())
	return diags
}

func fillTargetServer(c *client.TargetServer, d *schema.ResourceData) {
	isEnabled, ok := d.GetOk("is_enabled")
	if ok {
		c.IsEnabled = isEnabled.(bool)
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
		c.SSLInfo.ClientAuthEnabled = strconv.FormatBool(sslClientAuthEnabled.(bool))
	}
	sslIgnoreValidationErrors, ok := d.GetOk("ssl_ignore_validation_errors")
	if ok {
		c.SSLInfo.IgnoreValidationErrors = sslIgnoreValidationErrors.(bool)
	}
}

func resourceTargetServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.TargetServerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TargetServerPathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.TargetServer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("host", retVal.Host)
	d.Set("port", retVal.Port)
	d.Set("is_enabled", retVal.IsEnabled)
	if retVal.SSLInfo != nil {
		sslEnabled, _ := strconv.ParseBool(retVal.SSLInfo.Enabled)
		d.Set("ssl_enabled", sslEnabled)
		d.Set("ssl_keystore", retVal.SSLInfo.KeyStore)
		d.Set("ssl_keyalias", retVal.SSLInfo.KeyAlias)
		d.Set("ssl_truststore", retVal.SSLInfo.TrustStore)
		sslClientAuthEnabled, _ := strconv.ParseBool(retVal.SSLInfo.ClientAuthEnabled)
		d.Set("ssl_client_auth_enabled", sslClientAuthEnabled)
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

func resourceTargetServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.TargetServerDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upTargetServer := client.TargetServer{
		EnvironmentName: envName,
		Name:            name,
		Host:            d.Get("host").(string),
		Port:            d.Get("port").(int),
	}
	fillTargetServer(&upTargetServer, d)
	err := json.NewEncoder(&buf).Encode(upTargetServer)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.TargetServerPathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceTargetServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.TargetServerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.TargetServerPathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
