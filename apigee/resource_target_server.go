package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
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
			"ssl_common_name": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				DefaultFunc: getDefaultCommonName,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("COMMON_NAME", nil),
						},
						"wildcard_match": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"ssl_client_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_ignore_validation_errors": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTargetServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	var newTargetServer interface{}
	if c.IsGoogle() {
		newTS := client.GoogleTargetServer{
			EnvironmentName: d.Get("environment_name").(string),
			Name:            d.Get("name").(string),
			Host:            d.Get("host").(string),
			Port:            d.Get("port").(int),
		}
		fillGoogleTargetServer(&newTS, d)
		newTargetServer = &newTS
	} else {
		newTS := client.TargetServer{
			EnvironmentName: d.Get("environment_name").(string),
			Name:            d.Get("name").(string),
			Host:            d.Get("host").(string),
			Port:            d.Get("port").(int),
		}
		fillTargetServer(&newTS, d)
		newTargetServer = &newTS
	}
	err := json.NewEncoder(&buf).Encode(newTargetServer)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	var envName string
	if c.IsGoogle() {
		envName = newTargetServer.(*client.GoogleTargetServer).EnvironmentName
	} else {
		envName = newTargetServer.(*client.TargetServer).EnvironmentName
	}
	requestPath := fmt.Sprintf(client.TargetServerPath, c.Organization, envName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if c.IsGoogle() {
		d.SetId(newTargetServer.(*client.GoogleTargetServer).TargetServerEncodeId())
	} else {
		d.SetId(newTargetServer.(*client.TargetServer).TargetServerEncodeId())
	}
	return diags
}

func fillTargetServer(c *client.TargetServer, d *schema.ResourceData) {
	isEnabled, ok := d.GetOk("is_enabled")
	if ok {
		c.IsEnabled = isEnabled.(bool)
	}
	sslEnabled, ok := d.GetOk("ssl_enabled")
	if sslEnabled.(bool) {
		c.SSLInfo = &client.SSL{
			Enabled: strconv.FormatBool(true),
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
		sslCommonName, ok := d.GetOk("ssl_common_name.0")
		if ok {
			c.SSLInfo.CommonName = sslCommonName.(*client.SSLCommonName)

			value, ok := d.GetOk("ssl_common_name.0.value")
			if ok {
				c.SSLInfo.CommonName.Value = value.(string)
			} else {
				c.SSLInfo.CommonName.Value = c.Host
			}

			wildcard_match, ok := d.GetOk("ssl_common_name.0.wildcard_match")
			if ok {
				c.SSLInfo.CommonName.WildcardMatch = wildcard_match.(bool)
			}
		}
		protocols, ok := d.GetOk("protocols")
		if ok {
			set := protocols.(*schema.Set)
			c.SSLInfo.Protocols = convertSetToArray(set)
		}
		sslClientAuthEnabled, ok := d.GetOk("ssl_client_auth_enabled")
		c.SSLInfo.ClientAuthEnabled = strconv.FormatBool(sslClientAuthEnabled.(bool))
		sslIgnoreValidationErrors, ok := d.GetOk("ssl_ignore_validation_errors")
		c.SSLInfo.IgnoreValidationErrors = sslIgnoreValidationErrors.(bool)
	}
}

func fillGoogleTargetServer(c *client.GoogleTargetServer, d *schema.ResourceData) {
	isEnabled, ok := d.GetOk("is_enabled")
	if ok {
		c.IsEnabled = isEnabled.(bool)
	}
	sslEnabled, ok := d.GetOk("ssl_enabled")
	if sslEnabled.(bool) {
		c.SSLInfo = &client.GoogleSSL{
			Enabled: true,
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
		sslCommonName, ok := d.GetOk("ssl_common_name.0")
		if ok {
			c.SSLInfo.CommonName = sslCommonName.(*client.SSLCommonName)

			value, ok := d.GetOk("ssl_common_name.0.value")
			if ok {
				c.SSLInfo.CommonName.Value = value.(string)
			} else {
				c.SSLInfo.CommonName.Value = c.Host
			}

			wildcard_match, ok := d.GetOk("ssl_common_name.0.wildcard_match")
			if ok {
				c.SSLInfo.CommonName.WildcardMatch = wildcard_match.(bool)
			}
		}
		protocols, ok := d.GetOk("protocols")
		if ok {
			set := protocols.(*schema.Set)
			c.SSLInfo.Protocols = convertSetToArray(set)
		}
		sslClientAuthEnabled, ok := d.GetOk("ssl_client_auth_enabled")
		c.SSLInfo.ClientAuthEnabled = sslClientAuthEnabled.(bool)
		sslIgnoreValidationErrors, ok := d.GetOk("ssl_ignore_validation_errors")
		c.SSLInfo.IgnoreValidationErrors = sslIgnoreValidationErrors.(bool)
	}
}

func getDefaultCommonName() (interface{}, error) {
	if v := os.Getenv("COMMON_NAME"); v != "" {
		common_name := &client.SSLCommonName{Value: v}
		return common_name, nil
	}

	return nil, nil
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
	var retVal interface{}
	if c.IsGoogle() {
		retVal = &client.GoogleTargetServer{}
	} else {
		retVal = &client.TargetServer{}
	}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	var host, keyStore, keyAlias, trustStore string
	var port int
	var isEnabled, hasSSL, sslEnabled, clientAuthEnabled, ignoreValidationErrors bool
	var commonName *client.SSLCommonName
	var protocols []string
	if c.IsGoogle() {
		ts := retVal.(*client.GoogleTargetServer)
		host = ts.Host
		port = ts.Port
		isEnabled = ts.IsEnabled
		hasSSL = ts.SSLInfo != nil
		if hasSSL {
			keyStore = ts.SSLInfo.KeyStore
			keyAlias = ts.SSLInfo.KeyAlias
			trustStore = ts.SSLInfo.TrustStore
			commonName = ts.SSLInfo.CommonName
			sslEnabled = ts.SSLInfo.Enabled
			clientAuthEnabled = ts.SSLInfo.ClientAuthEnabled
			ignoreValidationErrors = ts.SSLInfo.IgnoreValidationErrors
			protocols = ts.SSLInfo.Protocols
		}
	} else {
		ts := retVal.(*client.TargetServer)
		host = ts.Host
		port = ts.Port
		isEnabled = ts.IsEnabled
		hasSSL = ts.SSLInfo != nil
		if hasSSL {
			keyStore = ts.SSLInfo.KeyStore
			keyAlias = ts.SSLInfo.KeyAlias
			trustStore = ts.SSLInfo.TrustStore
			commonName = ts.SSLInfo.CommonName
			sslEnabledBool, _ := strconv.ParseBool(ts.SSLInfo.Enabled)
			sslEnabled = sslEnabledBool
			clientAuthEnabledBool, _ := strconv.ParseBool(ts.SSLInfo.ClientAuthEnabled)
			clientAuthEnabled = clientAuthEnabledBool
			ignoreValidationErrors = ts.SSLInfo.IgnoreValidationErrors
			protocols = ts.SSLInfo.Protocols
		}
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("host", host)
	d.Set("port", port)
	d.Set("is_enabled", isEnabled)
	if hasSSL {
		d.Set("ssl_enabled", sslEnabled)
		d.Set("ssl_keystore", keyStore)
		d.Set("ssl_keyalias", keyAlias)
		d.Set("ssl_truststore", trustStore)
		d.Set("ssl_common_name.0", commonName)
		d.Set("ssl_client_auth_enabled", clientAuthEnabled)
		d.Set("ssl_ignore_validation_errors", ignoreValidationErrors)
		d.Set("protocols", protocols)
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
	var upTargetServer interface{}
	if c.IsGoogle() {
		upTS := client.GoogleTargetServer{
			EnvironmentName: envName,
			Name:            name,
			Host:            d.Get("host").(string),
			Port:            d.Get("port").(int),
		}
		fillGoogleTargetServer(&upTS, d)
		upTargetServer = &upTS
	} else {
		upTS := client.TargetServer{
			EnvironmentName: envName,
			Name:            name,
			Host:            d.Get("host").(string),
			Port:            d.Get("port").(int),
		}
		fillTargetServer(&upTS, d)
		upTargetServer = &upTS
	}
	err := json.NewEncoder(&buf).Encode(upTargetServer)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.TargetServerPathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
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
