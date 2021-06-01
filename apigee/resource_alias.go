package apigee

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
	"net/http"
	"net/url"
	"strconv"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"keystore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"keycertfile", "keycertjar", "pkcs12"}, false),
			},
			"file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_file_hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_file_hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ignore_expiry_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ignore_newline_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newAlias := client.Alias{
		EnvironmentName:         d.Get("environment_name").(string),
		KeystoreName:            d.Get("keystore_name").(string),
		Name:                    d.Get("name").(string),
		Format:                  d.Get("format").(string),
		IgnoreExpiryValidation:  d.Get("ignore_expiry_validation").(bool),
		IgnoreNewlineValidation: d.Get("ignore_newline_validation").(bool),
	}
	fd := map[string]client.FormData{}
	file, ok := d.GetOk("file")
	if ok {
		fd["file"] = client.FormData{Filename: file.(string)}
	}
	keyFile, ok := d.GetOk("key_file")
	if ok {
		fd["keyFile"] = client.FormData{Filename: keyFile.(string)}
	}
	certFile, ok := d.GetOk("cert_file")
	if ok {
		fd["certFile"] = client.FormData{Filename: certFile.(string)}
	}
	password, ok := d.GetOk("password")
	if ok {
		fd["password"] = client.FormData{Text: password.(string)}
	}
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(fd)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.AliasPath, c.Organization, newAlias.EnvironmentName, newAlias.KeystoreName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"alias":                   []string{newAlias.Name},
		"format":                  []string{newAlias.Format},
		"ignoreExpiryValidation":  []string{strconv.FormatBool(newAlias.IgnoreExpiryValidation)},
		"ignoreNewlineValidation": []string{strconv.FormatBool(newAlias.IgnoreNewlineValidation)},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newAlias.AliasEncodeId())
	return diags
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, keystoreName, name := client.AliasDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.AliasPathGet, c.Organization, envName, keystoreName, name)
	_, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("keystore_name", keystoreName)
	d.Set("name", name)
	//Use default values for flags
	d.Set("ignore_expiry_validation", false)
	d.Set("ignore_newline_validation", true)
	return diags
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, keystoreName, name := client.AliasDecodeId(d.Id())
	c := m.(*client.Client)
	//Only care about file changes
	if d.HasChanges("file", "file_hash") {
		file := d.Get("file").(string)
		//Turn filename into multi part buffer
		mp, buf, err := client.GetMultiPartBuffer(map[string]client.FormData{
			"file": client.FormData{Filename: file},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.AliasPathGet, c.Organization, envName, keystoreName, name)
		requestHeaders := http.Header{
			headers.ContentType: []string{mp.FormDataContentType()},
		}
		requestQuery := url.Values{
			"ignoreExpiryValidation": []string{strconv.FormatBool(d.Get("ignore_expiry_validation").(bool))},
		}
		_, err = c.HttpRequest(http.MethodPut, requestPath, requestQuery, requestHeaders, buf)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, keystoreName, name := client.AliasDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.AliasPathGet, c.Organization, envName, keystoreName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
