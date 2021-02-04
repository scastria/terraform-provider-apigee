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

func resourceCompanyDeveloper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCompanyDeveloperCreate,
		ReadContext:   resourceCompanyDeveloperRead,
		UpdateContext: resourceCompanyDeveloperUpdate,
		DeleteContext: resourceCompanyDeveloperDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"company_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"developer_email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCompanyDeveloperCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newCompanyDeveloper := client.CompanyDeveloper{
		CompanyName:    d.Get("company_name").(string),
		DeveloperEmail: d.Get("developer_email").(string),
	}
	fillCompanyDeveloper(&newCompanyDeveloper, d)
	newCompanyDevelopers := client.CompanyDeveloperList{
		Developers: []client.CompanyDeveloper{
			newCompanyDeveloper,
		},
	}
	err := json.NewEncoder(&buf).Encode(newCompanyDevelopers)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CompanyDeveloperPath, c.Organization, newCompanyDeveloper.CompanyName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newCompanyDeveloper.CompanyDeveloperEncodeId())
	return diags
}

func fillCompanyDeveloper(c *client.CompanyDeveloper, d *schema.ResourceData) {
	roleName, ok := d.GetOk("role_name")
	if ok {
		c.Role = roleName.(string)
	}
}

func resourceCompanyDeveloperRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, developerEmail := client.CompanyDeveloperDecodeId(d.Id())
	c := m.(*client.Client)
	//Must read all developers for this company and search for match
	requestPath := fmt.Sprintf(client.CompanyDeveloperPath, c.Organization, companyName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.CompanyDeveloperList{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Search for developer
	var foundCompanyDeveloper *client.CompanyDeveloper
	foundCompanyDeveloper = nil
	for _, cd := range retVal.Developers {
		if cd.DeveloperEmail == developerEmail {
			foundCompanyDeveloper = &cd
			break
		}
	}
	if foundCompanyDeveloper == nil {
		//Similar to 404 so don't return error
		d.SetId("")
		return diags
	}
	d.Set("company_name", companyName)
	d.Set("developer_email", developerEmail)
	d.Set("role_name", foundCompanyDeveloper.Role)
	return diags

}

func resourceCompanyDeveloperUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, developerEmail := client.CompanyDeveloperDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCompanyDeveloper := client.CompanyDeveloper{
		CompanyName:    companyName,
		DeveloperEmail: developerEmail,
	}
	fillCompanyDeveloper(&upCompanyDeveloper, d)
	upCompanyDevelopers := client.CompanyDeveloperList{
		Developers: []client.CompanyDeveloper{
			upCompanyDeveloper,
		},
	}
	err := json.NewEncoder(&buf).Encode(upCompanyDevelopers)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CompanyDeveloperPath, c.Organization, companyName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCompanyDeveloperDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	companyName, developerEmail := client.CompanyDeveloperDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CompanyDeveloperPathGet, c.Organization, companyName, developerEmail)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
