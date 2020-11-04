package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"email_id": {
				Type:     schema.TypeString,
				Required: true,
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
		return diag.FromErr(err)
	}
	body, err := c.HttpRequest("users", "POST", buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	d.SetId(retVal.EmailId)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	body, err := c.HttpRequest("users/"+d.Id(), "GET", bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	d.Set("first_name", retVal.FirstName)
	d.Set("last_name", retVal.LastName)
	d.SetId(retVal.EmailId)
	return diags

}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upUser := client.User{
		EmailId:   d.Get("email_id").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Password:  d.Get("password").(string),
	}
	err := json.NewEncoder(&buf).Encode(upUser)
	if err != nil {
		return diag.FromErr(err)
	}
	body, err := c.HttpRequest("users/"+upUser.EmailId, "PUT", buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	d.Set("first_name", retVal.FirstName)
	d.Set("last_name", retVal.LastName)
	d.SetId(retVal.EmailId)
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	body, err := c.HttpRequest("users/"+d.Id(), "DELETE", bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.User{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	return diags
}
