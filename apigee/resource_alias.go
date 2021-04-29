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
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"keystore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cert_file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(client.GetSupportedFormats(), false),
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"password_env_var_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func createAliasWithCert_file(c *client.Client, newAlias *client.Alias) (*client.Alias, error) {
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(interface{}(newAlias.CertFile).(string), "certFile")
	if err != nil {
		return nil, err
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
	body, err := c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		return nil, err
	}
	retVal := &client.Alias{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func createAliasWithMultipleFiles(c *client.Client, newAlias *client.Alias, password interface{}) (*client.Alias, error) {
	values := make(map[string]io.Reader)

	if len(newAlias.CertFile) != 0 {
		values["certFile"] = mustOpen(interface{}(newAlias.CertFile).(string))
	}

	if len(newAlias.File) != 0 {
		values["file"] = mustOpen(interface{}(newAlias.File).(string))
	}

	if len(newAlias.KeyFile) != 0 {
		values["keyFile"] = mustOpen(interface{}(newAlias.KeyFile).(string))
	}

	if password != nil {
		values["password"] = strings.NewReader(password.(string))
	}

	//Turn filenames/fields into multi part buffer
	mp, buf, err := client.GetMultiPartBufferForMultipleValues(values)
	if err != nil {
		return nil, err
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
	body, err := c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		return nil, err
	}
	retVal := &client.Alias{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}
func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	password := d.Get("password")
	c := m.(*client.Client)
	newAlias := client.Alias{
		EnvironmentName:         d.Get("environment_name").(string),
		Name:                    d.Get("name").(string),
		CertFile:                d.Get("cert_file").(string),
		File:                    d.Get("file").(string),
		KeyFile:                 d.Get("key_file").(string),
		KeystoreName:            d.Get("keystore_name").(string),
		Format:                  d.Get("format").(string),
		IgnoreExpiryValidation:  d.Get("ignore_expiry_validation").(bool),
		IgnoreNewlineValidation: d.Get("ignore_newline_validation").(bool),
	}

	if value, ok := d.GetOk("password"); ok {
		password = value.(string)
	}
	if value, ok := d.GetOk("password_env_var_name"); ok {
		password = os.Getenv(value.(string))
	}
	_, err := createAliasWithMultipleFiles(c, &newAlias, password.(string))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newAlias.AliasEncodeId())
	return diags
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	envName, keystoreName, name := client.AliasDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.AliasPathGet, c.Organization, envName, keystoreName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Alias{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("name", d.Id())
	//Retrieve the latest revision available as THE revision, assumes array is sorted
	return diags
}
func updateAliasWithCert_file(c *client.Client, newAlias *client.Alias) (*client.Alias, error) {
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(interface{}(newAlias.CertFile).(string), "certFile")
	if err != nil {
		return nil, err
	}
	requestPath := fmt.Sprintf(client.AliasPathUpdate, c.Organization, newAlias.EnvironmentName, newAlias.KeystoreName, newAlias.Name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"ignoreExpiryValidation": []string{strconv.FormatBool(newAlias.IgnoreExpiryValidation)},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		return nil, err
	}
	retVal := &client.Alias{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	upAlias := client.Alias{
		EnvironmentName:         d.Get("environment_name").(string),
		Name:                    d.Id(),
		CertFile:                d.Get("cert_file").(string),
		KeystoreName:            d.Get("keystore_name").(string),
		Format:                  d.Get("format").(string),
		IgnoreExpiryValidation:  d.Get("ignore_expiry_validation").(bool),
		IgnoreNewlineValidation: d.Get("ignore_newline_validation").(bool),
	}

	_, err := updateAliasWithCert_file(c, &upAlias)

	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	envName, keystoreName, name := client.AliasDecodeId(d.Id())
	requestPath := fmt.Sprintf(client.AliasPathGet, c.Organization, envName , keystoreName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
