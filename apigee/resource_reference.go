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
)

func resourceReference() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReferenceCreate,
		ReadContext:   resourceReferenceRead,
		UpdateContext: resourceReferenceUpdate,
		DeleteContext: resourceReferenceDelete,
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
			"refers": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"KeyStore", "TrustStore"}, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceReferenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newReference := client.Reference{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
		Refers:          d.Get("refers").(string),
		ResourceType:    d.Get("resource_type").(string),
	}
	err := json.NewEncoder(&buf).Encode(newReference)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ReferencePath, c.Organization, newReference.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newReference.ReferenceEncodeId())
	return diags
}

func resourceReferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.ReferenceDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ReferencePathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Reference{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("refers", retVal.Refers)
	d.Set("resource_type", retVal.ResourceType)

	return diags
}

func resourceReferenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.ReferenceDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upReference := client.Reference{
		EnvironmentName: envName,
		Name:            name,
	}
	fillReference(&upReference, d)
	err := json.NewEncoder(&buf).Encode(upReference)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ReferencePathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func fillReference(c *client.Reference, d *schema.ResourceData) {
	refers, ok := d.GetOk("refers")
	if ok {
		c.Refers = refers.(string)
	}
	resourceType, ok := d.GetOk("resource_type")
	if ok {
		c.ResourceType = resourceType.(string)
	}
}

func resourceReferenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.ReferenceDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ReferencePathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
