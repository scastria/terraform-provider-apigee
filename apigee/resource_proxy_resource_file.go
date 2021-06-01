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

func resourceProxyResourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyResourceFileCreate,
		ReadContext:   resourceProxyResourceFileRead,
		UpdateContext: resourceProxyResourceFileUpdate,
		DeleteContext: resourceProxyResourceFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"proxy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
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

func resourceProxyResourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newProxyResourceFile := client.ResourceFile{
		ProxyName: d.Get("proxy_name").(string),
		Revision:  d.Get("revision").(int),
		Type:      d.Get("type").(string),
		Name:      d.Get("name").(string),
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
	requestPath := fmt.Sprintf(client.ProxyResourceFilePath, c.Organization, newProxyResourceFile.ProxyName, newProxyResourceFile.Revision)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"type": []string{newProxyResourceFile.Type},
		"name": []string{newProxyResourceFile.Name},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newProxyResourceFile.ProxyResourceFileEncodeId())
	return diags
}

func resourceProxyResourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	proxyName, rev, rtype, name := client.ProxyResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	//Reading specific file returns actual contents of file so read all files of type and search for name instead
	requestPath := fmt.Sprintf(client.ProxyResourceFilePathOfType, c.Organization, proxyName, rev, rtype)
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
	d.Set("proxy_name", proxyName)
	d.Set("revision", rev)
	d.Set("type", rtype)
	d.Set("name", name)
	return diags
}

func resourceProxyResourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	proxyName, rev, rtype, name := client.ProxyResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	file := d.Get("file").(string)
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(map[string]client.FormData{
		"file": client.FormData{Filename: file},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ProxyResourceFilePathGet, c.Organization, proxyName, rev, rtype, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceProxyResourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	proxyName, rev, rtype, name := client.ProxyResourceFileDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyResourceFilePathGet, c.Organization, proxyName, rev, rtype, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
