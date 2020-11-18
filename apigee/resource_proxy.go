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
	"strconv"
)

func resourceProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyCreate,
		ReadContext:   resourceProxyRead,
		UpdateContext: resourceProxyUpdate,
		DeleteContext: resourceProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bundle": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bundle_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func importRevision(c *client.Client, name string, bundle string) (*client.ProxyRevision, error) {
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(bundle, "bundle")
	if err != nil {
		return nil, err
	}
	requestPath := fmt.Sprintf(client.ProxyPath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"action": []string{"import"},
		"name":   []string{name},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		return nil, err
	}
	retVal := &client.ProxyRevision{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func resourceProxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Get("name").(string)
	bundle := d.Get("bundle").(string)
	retVal, err := importRevision(c, name, bundle)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(name)
	revision, _ := strconv.Atoi(retVal.Revision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyPathGet, c.Organization, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Proxy{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("name", d.Id())
	//Retrieve the latest revision available as THE revision
	lastRevision := retVal.Revisions[len(retVal.Revisions)-1]
	revision, _ := strconv.Atoi(lastRevision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Get("name").(string)
	bundle := d.Get("bundle").(string)
	retVal, err := importRevision(c, name, bundle)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	revision, _ := strconv.Atoi(retVal.Revision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyPathGet, c.Organization, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
