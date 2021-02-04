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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email_id": {
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
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newUser := client.User{
		EmailId:   d.Get("email_id").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Password:  d.Get("password").(string),
	}
	err := json.NewEncoder(&buf).Encode(newUser)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.UserPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newUser.EmailId)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("email_id", d.Id())
	d.Set("first_name", retVal.FirstName)
	d.Set("last_name", retVal.LastName)
	//Cannot set password since it is not returned by API to keep secret
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	//Do not use id since that can change on update
	upUser := client.User{
		EmailId:   d.Get("email_id").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
	}
	if d.HasChange("password") {
		upUser.Password = d.Get("password").(string)
	}
	err := json.NewEncoder(&buf).Encode(upUser)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	//EmailId can be changed which changes the id
	d.SetId(upUser.EmailId)
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserPathGet, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
