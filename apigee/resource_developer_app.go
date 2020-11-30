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
	"regexp"
)

func resourceDeveloperApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeveloperAppCreate,
		ReadContext:   resourceDeveloperAppRead,
		UpdateContext: resourceDeveloperAppUpdate,
		DeleteContext: resourceDeveloperAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"developer_email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
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

func resourceDeveloperAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newDeveloperApp := client.DeveloperApp{
		DeveloperEmail: d.Get("developer_email").(string),
		Name:           d.Get("name").(string),
	}
	fillDeveloperApp(&newDeveloperApp, d)
	err := json.NewEncoder(&buf).Encode(newDeveloperApp)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperAppPath, c.Organization, newDeveloperApp.DeveloperEmail)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Set id before decoding to get generated key for deletion
	d.SetId(newDeveloperApp.DeveloperAppEncodeId())
	retVal := &client.DeveloperApp{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		//Don't clear id since app was created
		return diag.FromErr(err)
	}
	//Delete generated keys so that user is in control of keys via Terraform
	for _, key := range retVal.Credentials {
		requestPath = fmt.Sprintf(client.DeveloperAppPathGeneratedKey, c.Organization, newDeveloperApp.DeveloperEmail, newDeveloperApp.Name, key.ConsumerKey)
		_, err = c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			//Don't clear id since app was created
			return diag.FromErr(err)
		}
	}
	return diags
}

func fillDeveloperApp(c *client.DeveloperApp, d *schema.ResourceData) {
	callback, ok := d.GetOk("callback_url")
	if ok {
		c.CallbackURL = callback.(string)
	}
	a, ok := d.GetOk("attributes")
	if ok {
		attributes := a.(map[string]interface{})
		for name, value := range attributes {
			c.Attributes = append(c.Attributes, client.DeveloperAppAttribute{
				Name:  name,
				Value: value.(string),
			})
		}
	}
}

func resourceDeveloperAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, name := client.DeveloperAppDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperAppPathGet, c.Organization, developerEmail, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.DeveloperApp{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("developer_email", developerEmail)
	d.Set("name", name)
	d.Set("callback_url", retVal.CallbackURL)
	atts := map[string]string{}
	for _, e := range retVal.Attributes {
		atts[e.Name] = e.Value
	}
	d.Set("attributes", atts)
	return diags
}

func resourceDeveloperAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, name := client.DeveloperAppDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upDeveloperApp := client.DeveloperApp{
		DeveloperEmail: developerEmail,
		Name:           name,
	}
	fillDeveloperApp(&upDeveloperApp, d)
	err := json.NewEncoder(&buf).Encode(upDeveloperApp)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperAppPathGet, c.Organization, developerEmail, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDeveloperAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, name := client.DeveloperAppDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperAppPathGet, c.Organization, developerEmail, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
