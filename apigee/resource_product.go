package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
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
			"operation_config_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) ([]string, []error) {
					value := val.(string)
					if value != "proxy" && value != "remoteservice" {
						return []string{}, []error{fmt.Errorf("Invalid value for %s, must be either 'proxy' or 'remoteservice'", key)}
					}
					return nil, nil
				},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"operation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"methods": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
						"attributes": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceProductCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		headers.ContentType: []string{client.ApplicationJson},
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
	ops, ok := d.GetOk("operation")
	if ok {
		oct, ok2 := d.GetOk("operation_config_type")
		if ok2 {
			fillOperationsConfig(c, ops, oct.(string))
		} else {
			fillOperationsConfig(c, ops, "proxy")
		}
	}
}

func fillOperationsConfig(c *client.Product, ops interface{}, oct string) {
	operations := ops.([]interface{})
	c.OperationGroup.OperationConfigs = make([]client.OperationConfigs, len(operations))
	c.OperationGroup.OperationConfigType = oct
	for i, op := range operations {
		item := op.(map[string]interface{})
		q := client.Quota{}

		quota, ok := item["quota"].(int)
		if ok && quota > 0 {
			q.Limit = strconv.Itoa(quota)
			q.Interval = strconv.Itoa(item["quota_interval"].(int))
			q.TimeUnit = item["quota_time_unit"].(string)
		}

		attributes := item["attributes"].(map[string]interface{})
		var attribs []client.Attribute
		for name, value := range attributes {
			attribs = append(attribs, client.Attribute{
				Name:  name,
				Value: value.(string),
			})
		}

		c.OperationGroup.OperationConfigs[i] = client.OperationConfigs{
			ApiSource: item["api_source"].(string),
			Operations: []client.Operation{{
				Resource: item["path"].(string),
				Methods:  convertSetToArray(item["methods"].(*schema.Set)),
			}},
			Quota:      q,
			Attributes: attribs,
		}
	}
}

func readOperationsConfig(c client.OperationGroup) []interface{} {
	operations := make([]interface{}, len(c.OperationConfigs))
	for i, config := range c.OperationConfigs {
		atts := map[string]string{}
		for _, e := range config.Attributes {
			atts[e.Name] = e.Value
		}
		operations[i] = map[string]interface{}{
			"api_source":      config.ApiSource,
			"path":            config.Operations[0].Resource,
			"methods":         config.Operations[0].Methods,
			"quota":           config.Quota.Limit,
			"quota_interval":  config.Quota.Interval,
			"quota_time_unit": config.Quota.TimeUnit,
			"attributes":      atts,
		}
	}
	return operations
}

func resourceProductRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	d.Set("operations", readOperationsConfig(retVal.OperationGroup))
	d.Set("operaton_config_type", retVal.OperationGroup.OperationConfigType)
	return diags
}

func resourceProductUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceProductDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
