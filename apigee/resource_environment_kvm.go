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

func resourceEnvironmentKVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentKVMCreate,
		ReadContext:   resourceEnvironmentKVMRead,
		UpdateContext: resourceEnvironmentKVMUpdate,
		DeleteContext: resourceEnvironmentKVMDelete,
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
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"entry": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					Sensitive: true,
				},
			},
		},
		CustomizeDiff: resourceEnvironmentKVMCustomDiff,
	}
}

func resourceEnvironmentKVMCustomDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	//A KVM cannot be decrypted so ForceNew if going from true to false
	if !diff.HasChange("encrypted") {
		return nil
	}
	o, n := diff.GetChange("encrypted")
	if o.(bool) && !n.(bool) {
		diff.ForceNew("encrypted")
	}
	return nil
}

func resourceEnvironmentKVMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newEnvironmentKVM := client.EnvironmentKVM{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
	}
	fillEnvironmentKVM(&newEnvironmentKVM, d)
	err := json.NewEncoder(&buf).Encode(newEnvironmentKVM)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentKVMPath, c.Organization, newEnvironmentKVM.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newEnvironmentKVM.EnvironmentKVMEncodeId())
	return diags
}

func fillEnvironmentKVM(c *client.EnvironmentKVM, d *schema.ResourceData) {
	encrypted, ok := d.GetOk("encrypted")
	if ok {
		c.Encrypted = encrypted.(bool)
	}
	e, ok := d.GetOk("entry")
	if ok {
		entries := e.(map[string]interface{})
		for name, value := range entries {
			c.Entries = append(c.Entries, client.KVMEntry{
				Name:  name,
				Value: value.(string),
			})
		}
	}
}

func resourceEnvironmentKVMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.EnvironmentKVMDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnvironmentKVMPathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.EnvironmentKVM{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	entries := map[string]string{}
	for _, e := range retVal.Entries {
		entries[e.Name] = e.Value
	}
	d.Set("entry", entries)
	return diags
}

func resourceEnvironmentKVMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.EnvironmentKVMDecodeId(d.Id())
	c := m.(*client.Client)
	//Check for removal of entries
	if d.HasChange("entry") {
		o, n := d.GetChange("entry")
		old := o.(map[string]interface{})
		new := n.(map[string]interface{})
		for oldKey, _ := range old {
			_, newHasKey := new[oldKey]
			if newHasKey {
				continue
			}
			//Delete entry
			requestPath := fmt.Sprintf(client.EnvironmentKVMPathGetEntry, c.Organization, envName, name, oldKey)
			_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	buf := bytes.Buffer{}
	upEnvironmentKVM := client.EnvironmentKVM{
		EnvironmentName: envName,
		Name:            name,
	}
	fillEnvironmentKVM(&upEnvironmentKVM, d)
	err := json.NewEncoder(&buf).Encode(upEnvironmentKVM)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentKVMPathGet, c.Organization, envName, name)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceEnvironmentKVMDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.EnvironmentKVMDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnvironmentKVMPathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
