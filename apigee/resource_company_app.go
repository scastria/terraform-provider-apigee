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
	"mime"
	"net/http"
)

func resourceCompanyApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCompanyAppCreate,
		ReadContext:   resourceCompanyAppRead,
		UpdateContext: resourceCompanyAppUpdate,
		DeleteContext: resourceCompanyAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"callback_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCompanyAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newCompanyApp := client.App{
		CompanyName: d.Get("company_name").(string),
		Name:        d.Get("name").(string),
	}
	fillCompanyApp(&newCompanyApp, d)
	err := json.NewEncoder(&buf).Encode(newCompanyApp)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CompanyAppPath, c.Organization, newCompanyApp.CompanyName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Set id before decoding to get generated key for deletion
	d.SetId(newCompanyApp.CompanyAppEncodeId())
	retVal := &client.App{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		//Don't clear id since app was created
		return diag.FromErr(err)
	}
	//Delete generated keys so that user is in control of keys via Terraform
	for _, key := range retVal.Credentials {
		requestPath = fmt.Sprintf(client.CompanyAppPathGeneratedKey, c.Organization, newCompanyApp.CompanyName, newCompanyApp.Name, key.ConsumerKey)
		_, err = c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			//Don't clear id since app was created
			return diag.FromErr(err)
		}
	}
	return diags
}

func fillCompanyApp(c *client.App, d *schema.ResourceData) {
	callback, ok := d.GetOk("callback_url")
	if ok {
		c.CallbackURL = callback.(string)
	}
	a, ok := d.GetOk("attributes")
	if ok {
		attributes := a.(map[string]interface{})
		for name, value := range attributes {
			c.Attributes = append(c.Attributes, client.Attribute{
				Name:  name,
				Value: value.(string),
			})
		}
	}
}

func resourceCompanyAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, name := client.AppDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CompanyAppPathGet, c.Organization, companyName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.App{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("company_name", companyName)
	d.Set("name", name)
	d.Set("callback_url", retVal.CallbackURL)
	atts := map[string]string{}
	for _, e := range retVal.Attributes {
		atts[e.Name] = e.Value
	}
	d.Set("attributes", atts)
	return diags
}

func resourceCompanyAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, name := client.AppDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCompanyApp := client.App{
		CompanyName: companyName,
		Name:        name,
	}
	fillCompanyApp(&upCompanyApp, d)
	err := json.NewEncoder(&buf).Encode(upCompanyApp)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CompanyAppPathGet, c.Organization, companyName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCompanyAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, name := client.AppDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CompanyAppPathGet, c.Organization, companyName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
