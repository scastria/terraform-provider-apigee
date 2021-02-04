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

func resourceRolePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRolePermissionCreate,
		ReadContext:   resourceRolePermissionRead,
		DeleteContext: resourceRolePermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"get", "put", "delete"}, false),
				},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRolePermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	permSet := d.Get("permissions").(*schema.Set)
	permList := convertSetToArray(permSet)
	if len(permList) == 0 {
		d.SetId("")
		return diag.Errorf("permissions must contain at least 1 value")
	}
	newRolePermission := client.RolePermission{
		RoleName:    d.Get("role_name").(string),
		Path:        d.Get("path").(string),
		Permissions: permList,
	}
	err := json.NewEncoder(&buf).Encode(newRolePermission)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RolePermissionPath, c.Organization, newRolePermission.RoleName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newRolePermission.RolePermissionEncodeId())
	return diags
}

func resourceRolePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	roleName, path := client.RolePermissionDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RolePermissionPath, c.Organization, roleName)
	requestQuery := url.Values{
		"path": []string{path},
	}
	body, err := c.HttpRequest(http.MethodGet, requestPath, requestQuery, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.RolePermission{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("role_name", roleName)
	d.Set("path", path)
	d.Set("permissions", retVal.Permissions)
	return diags
}

func resourceRolePermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	roleName, path := client.RolePermissionDecodeId(d.Id())
	c := m.(*client.Client)
	requestQuery := url.Values{
		"path": []string{path},
	}
	//Remove each perm
	for _, p := range []string{"get", "put", "delete"} {
		requestPath := fmt.Sprintf(client.RolePermissionPathGet, c.Organization, roleName, p)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, requestQuery, nil, &bytes.Buffer{})
		if err != nil {
			re := err.(*client.RequestError)
			if re.StatusCode != http.StatusNotFound {
				return diag.FromErr(err)
			}
		}
	}
	d.SetId("")
	return diags
}
