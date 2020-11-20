package apigee

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
	"net/http"
	"net/url"
	"regexp"
)

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"email_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newUserRole := client.UserRole{
		EmailId:  d.Get("email_id").(string),
		RoleName: d.Get("role_name").(string),
	}
	requestPath := fmt.Sprintf(client.UserRolePath, c.Organization, newUserRole.RoleName)
	requestForm := url.Values{
		"id": []string{newUserRole.EmailId},
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.FormEncoded},
	}
	_, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, bytes.NewBufferString(requestForm.Encode()))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newUserRole.UserRoleEncodeId())
	return diags
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	emailId, roleName := client.UserRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserRolePathGet, c.Organization, roleName, emailId)
	_, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("email_id", emailId)
	d.Set("role_name", roleName)
	return diags

}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	emailId, roleName := client.UserRoleDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.UserRolePathGet, c.Organization, roleName, emailId)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
