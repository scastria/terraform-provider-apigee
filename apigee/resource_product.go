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
	"strconv"
)

func resourceProduct() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProductCreate,
		ReadContext:   resourceProductRead,
		UpdateContext: resourceProductUpdate,
		DeleteContext: resourceProductDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_approval_type": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quota": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				RequiredWith: []string{"quota_interval", "quota_time_unit"},
			},
			"quota_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				RequiredWith: []string{"quota", "quota_time_unit"},
			},
			"quota_time_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"minute", "hour", "day", "month"}, false),
				RequiredWith: []string{"quota", "quota_interval"},
			},
			"api_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environments": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"proxies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceProductCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	autoApprovalType := d.Get("auto_approval_type").(bool)
	approvalType := ""
	if autoApprovalType {
		approvalType = client.AutoApprovalType
	} else {
		approvalType = client.ManualApprovalType
	}
	newProduct := client.Product{
		Name:         d.Get("name").(string),
		DisplayName:  d.Get("display_name").(string),
		ApprovalType: approvalType,
	}
	fillProduct(&newProduct, d)
	err := json.NewEncoder(&buf).Encode(newProduct)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ProductPath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newProduct.Name)
	return diags
}

func fillProduct(c *client.Product, d *schema.ResourceData) {
	desc, ok := d.GetOk("description")
	if ok {
		c.Description = desc.(string)
	}
	quota, ok := d.GetOk("quota")
	if ok {
		c.Quota = strconv.Itoa(quota.(int))
	}
	quotaInterval, ok := d.GetOk("quota_interval")
	if ok {
		c.QuotaInterval = strconv.Itoa(quotaInterval.(int))
	}
	quotaTimeUnit, ok := d.GetOk("quota_time_unit")
	if ok {
		c.QuotaTimeUnit = quotaTimeUnit.(string)
	}
	apiResources, ok := d.GetOk("api_resources")
	if ok {
		set := apiResources.(*schema.Set)
		c.APIResources = convertSetToArray(set)
	}
	envs, ok := d.GetOk("environments")
	if ok {
		set := envs.(*schema.Set)
		c.Environments = convertSetToArray(set)
	}
	proxies, ok := d.GetOk("proxies")
	if ok {
		set := proxies.(*schema.Set)
		c.Proxies = convertSetToArray(set)
	}
	scopes, ok := d.GetOk("scopes")
	if ok {
		set := scopes.(*schema.Set)
		c.Scopes = convertSetToArray(set)
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

func resourceProductRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProductPathGet, c.Organization, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Product{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("name", d.Id())
	d.Set("display_name", retVal.DisplayName)
	d.Set("auto_approval_type", retVal.ApprovalType == "auto")
	d.Set("description", retVal.Description)
	quota := retVal.Quota
	if quota != "" {
		quotaInt, _ := strconv.Atoi(quota)
		d.Set("quota", quotaInt)
	}
	quotaInterval := retVal.QuotaInterval
	if quotaInterval != "" {
		quotaIntervalInt, _ := strconv.Atoi(quotaInterval)
		d.Set("quota_interval", quotaIntervalInt)
	}
	d.Set("quota_time_unit", retVal.QuotaTimeUnit)
	d.Set("api_resources", retVal.APIResources)
	d.Set("environments", retVal.Environments)
	d.Set("proxies", retVal.Proxies)
	d.Set("scopes", retVal.Scopes)
	atts := map[string]string{}
	for _, e := range retVal.Attributes {
		atts[e.Name] = e.Value
	}
	d.Set("attributes", atts)
	return diags
}

func resourceProductUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	autoApprovalType := d.Get("auto_approval_type").(bool)
	approvalType := ""
	if autoApprovalType {
		approvalType = client.AutoApprovalType
	} else {
		approvalType = client.ManualApprovalType
	}
	upProduct := client.Product{
		Name:         d.Id(),
		DisplayName:  d.Get("display_name").(string),
		ApprovalType: approvalType,
	}
	fillProduct(&upProduct, d)
	err := json.NewEncoder(&buf).Encode(upProduct)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ProductPathGet, c.Organization, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceProductDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProductPathGet, c.Organization, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
