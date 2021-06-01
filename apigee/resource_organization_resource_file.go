package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
	"net/http"
	"net/url"
)

func resourceOrganizationResourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationResourceFileCreate,
		ReadContext:   resourceOrganizationResourceFileRead,
		UpdateContext: resourceOrganizationResourceFileUpdate,
		DeleteContext: resourceOrganizationResourceFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceOrganizationResourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newOrganizationResourceFile := client.ResourceFile{
		Type: d.Get("type").(string),
		Name: d.Get("name").(string),
	}
	file := d.Get("file").(string)
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(map[string]client.FormData{
		"file": client.FormData{Filename: file},
	})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.OrganizationResourceFilePath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"type": []string{newOrganizationResourceFile.Type},
		"name": []string{newOrganizationResourceFile.Name},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newOrganizationResourceFile.OrganizationResourceFileEncodeId())
	return diags
}

func resourceOrganizationResourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	rtype, name := client.OrganizationResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	//Reading specific file returns actual contents of file so read all files of type and search for name instead
	requestPath := fmt.Sprintf(client.OrganizationResourceFilePathOfType, c.Organization, rtype)
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
	d.Set("type", rtype)
	d.Set("name", name)
	return diags
}

func resourceOrganizationResourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	rtype, name := client.OrganizationResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	file := d.Get("file").(string)
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(map[string]client.FormData{
		"file": client.FormData{Filename: file},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.OrganizationResourceFilePathGet, c.Organization, rtype, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceOrganizationResourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	rtype, name := client.OrganizationResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.OrganizationResourceFilePathGet, c.Organization, rtype, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
