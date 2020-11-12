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

func resourceCache() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCacheCreate,
		ReadContext:   resourceCacheRead,
		UpdateContext: resourceCacheUpdate,
		DeleteContext: resourceCacheDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiry_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout_in_sec": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_of_day": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expiry_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"overflow_to_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"skip_cache_if_element_size_in_kb_exceeds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceCacheCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newCache := client.Cache{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
	}
	err := json.NewEncoder(&buf).Encode(newCache)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CachePath, c.Organization, newCache.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newCache.CacheEncodeId())
	return diags
}

func resourceCacheRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.CacheDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CachePathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Cache{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("description", retVal.Description)
	expirySettings := map[string]interface{}{}
	if retVal.ExpirySettings.TimeoutInSec != nil {
		expirySettings["timeout_in_sec"] = retVal.ExpirySettings.TimeoutInSec.Value
	} else if retVal.ExpirySettings.TimeOfDay != nil {
		expirySettings["time_of_day"] = retVal.ExpirySettings.TimeOfDay.Value
	} else if retVal.ExpirySettings.ExpiryDate != nil {
		expirySettings["expiry_date"] = retVal.ExpirySettings.ExpiryDate.Value
	}
	d.Set("expiry_settings", []interface{}{expirySettings})
	d.Set("overflow_to_disk", retVal.OverflowToDisk)
	return diags
}

func resourceCacheUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.CacheDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upCache := client.Cache{
		EnvironmentName: envName,
		Name:            name,
		Description:     d.Get("description").(string),
	}
	expirySettingsList := d.Get("expiry_settings").([]interface{})
	if len(expirySettingsList) == 1 {
		upCache.Description = "Bad"
		expirySettings := expirySettingsList[0].(map[string]interface{})
		timeoutInSec, timeoutInSecOk := expirySettings["timeout_in_sec"]
		timeOfDay, timeOfDayOk := expirySettings["time_of_day"]
		expiryDate, expiryDateOk := expirySettings["expiry_date"]
		if (timeoutInSecOk) && (timeoutInSec.(string) != "") {
			upCache.ExpirySettings.TimeoutInSec = &client.ExpiryValue{
				Value: timeoutInSec.(string),
			}
		} else if (timeOfDayOk) && (timeOfDay.(string) != "") {
			upCache.ExpirySettings.TimeOfDay = &client.ExpiryValue{
				Value: timeOfDay.(string),
			}
		} else if (expiryDateOk) && (expiryDate.(string) != "") {
			upCache.ExpirySettings.ExpiryDate = &client.ExpiryValue{
				Value: expiryDate.(string),
			}
		}
	}
	err := json.NewEncoder(&buf).Encode(upCache)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.CachePathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceCacheDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.CacheDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CachePathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
