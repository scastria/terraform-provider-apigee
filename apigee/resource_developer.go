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
	"regexp"
)

func resourceDeveloper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeveloperCreate,
		ReadContext:   resourceDeveloperRead,
		UpdateContext: resourceDeveloperUpdate,
		DeleteContext: resourceDeveloperDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceDeveloperCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newDeveloper := client.Developer{
		Email:     d.Get("email").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		UserName:  d.Get("user_name").(string),
	}
	fillDeveloper(&newDeveloper, d)
	err := json.NewEncoder(&buf).Encode(newDeveloper)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperPath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newDeveloper.Email)
	return diags
}

func fillDeveloper(c *client.Developer, d *schema.ResourceData) {
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

func resourceDeveloperRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperPathGet, c.Organization, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Developer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("email", d.Id())
	d.Set("first_name", retVal.FirstName)
	d.Set("last_name", retVal.LastName)
	d.Set("user_name", retVal.UserName)
	atts := map[string]string{}
	for _, e := range retVal.Attributes {
		atts[e.Name] = e.Value
	}
	d.Set("attributes", atts)
	return diags
}

func resourceDeveloperUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	//Do not use id since that can change on update
	upDeveloper := client.Developer{
		Email:     d.Get("email").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		UserName:  d.Get("user_name").(string),
	}
	fillDeveloper(&upDeveloper, d)
	err := json.NewEncoder(&buf).Encode(upDeveloper)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperPathGet, c.Organization, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	//Email can be changed which changes the id
	d.SetId(upDeveloper.Email)
	return diags
}

func resourceDeveloperDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperPathGet, c.Organization, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
