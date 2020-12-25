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
)

func resourceProxyPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyPolicyCreate,
		ReadContext:   resourceProxyPolicyRead,
		DeleteContext: resourceProxyPolicyDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_hash": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProxyPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newProxyPolicy := client.ProxyPolicy{
		ProxyName: d.Get("proxy_name").(string),
		Revision:  d.Get("revision").(int),
		Name:      d.Get("name").(string),
	}
	file := d.Get("file").(string)
	//Turn filename into buffer
	buf, err := client.GetBuffer(file)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ProxyPolicyPath, c.Organization, newProxyPolicy.ProxyName, newProxyPolicy.Revision)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".xml")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newProxyPolicy.ProxyPolicyEncodeId())
	return diags
}

func resourceProxyPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	proxyName, rev, name := client.ProxyPolicyDecodeId(d.Id())
	c := m.(*client.Client)
	//Reading specific file returns actual contents of file so read all policies and search for name instead
	requestPath := fmt.Sprintf(client.ProxyPolicyPath, c.Organization, proxyName, rev)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(body)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := []string{}
	err = json.Unmarshal(respBody.Bytes(), &retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Look for existence of name
	found := false
	for _, file := range retVal {
		if file == name {
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
	d.Set("name", name)
	return diags
}

func resourceProxyPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	proxyName, rev, name := client.ProxyPolicyDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyPolicyPathGet, c.Organization, proxyName, rev, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
