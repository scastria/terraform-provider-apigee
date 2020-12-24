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
	"net/http"
	"net/url"
)

func resourceEnvironmentResourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentResourceFileCreate,
		ReadContext:   resourceEnvironmentResourceFileRead,
		UpdateContext: resourceEnvironmentResourceFileUpdate,
		DeleteContext: resourceEnvironmentResourceFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"java", "js", "jsc", "hosted", "node", "py", "wsdl", "xsd", "xsl"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceEnvironmentResourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newEnvironmentResourceFile := client.ResourceFile{
		EnvironmentName: d.Get("environment_name").(string),
		Type:            d.Get("type").(string),
		Name:            d.Get("name").(string),
	}
	file := d.Get("file").(string)
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(file, "file")
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentResourceFilePath, c.Organization, newEnvironmentResourceFile.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"type": []string{newEnvironmentResourceFile.Type},
		"name": []string{newEnvironmentResourceFile.Name},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newEnvironmentResourceFile.EnvironmentResourceFileEncodeId())
	return diags
}

func resourceEnvironmentResourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, rtype, name := client.EnvironmentResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	//Reading specific file returns actual contents of file so read all files of type and search for name instead
	requestPath := fmt.Sprintf(client.EnvironmentResourceFilePathOfType, c.Organization, envName, rtype)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ResourceFilesOfType{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Look for existence of name
	found := false
	for _, file := range retVal.Files {
		if file.Name == name {
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
		return diags
	}
	d.Set("environment_name", envName)
	d.Set("type", rtype)
	d.Set("name", name)
	return diags
}

func resourceEnvironmentResourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, rtype, name := client.EnvironmentResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	file := d.Get("file").(string)
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(file, "file")
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentResourceFilePathGet, c.Organization, envName, rtype, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceEnvironmentResourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, rtype, name := client.EnvironmentResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnvironmentResourceFilePathGet, c.Organization, envName, rtype, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
