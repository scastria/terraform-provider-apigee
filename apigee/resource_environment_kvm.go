package apigee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-apigee/apigee/client"
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
				ForceNew: true,
			},
			"entry": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"sensitive_entry"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sensitive_entry": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"entry"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceEnvironmentKVMCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newEnvironmentKVM := client.KVM{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
	}
	fillEnvironmentKVM(&newEnvironmentKVM, c.IsPublic(), d)
	err := json.NewEncoder(&buf).Encode(newEnvironmentKVM)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.EnvironmentKVMPath, c.Organization, newEnvironmentKVM.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newEnvironmentKVM.EnvironmentKVMEncodeId())

	// Create individual KVM's on Apigee public
	if c.IsPublic() {
		e, ok := d.GetOk("sensitive_entry")
		if ok {
			entries := e.(map[string]interface{})
			for name, value := range entries {
				buf = bytes.Buffer{}
				err = json.NewEncoder(&buf).Encode(client.Attribute{
					Name:  name,
					Value: value.(string),
				})
				if err != nil {
					return diag.FromErr(err)
				}
				requestPath = fmt.Sprintf(client.EnvironmentKVMPathEntries, c.Organization, newEnvironmentKVM.EnvironmentName, newEnvironmentKVM.Name)
				_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
				if err != nil {
					return diag.FromErr(err)
				}

				newEnvironmentKVM.Entries = append(newEnvironmentKVM.Entries, client.Attribute{
					Name:  name,
					Value: value.(string),
				})
			}
		}
	}

	return diags
}

func fillEnvironmentKVM(c *client.KVM, isPublic bool, d *schema.ResourceData) {
	encrypted, ok := d.GetOk("encrypted")
	if ok {
		c.Encrypted = encrypted.(bool)
	}
	var e interface{}
	if c.Encrypted {
		e, ok = d.GetOk("sensitive_entry")
	} else {
		e, ok = d.GetOk("entry")
	}
	if ok && !isPublic {
		entries := e.(map[string]interface{})
		for name, value := range entries {
			c.Entries = append(c.Entries, client.Attribute{
				Name:  name,
				Value: value.(string),
			})
		}
	}
}

func resourceEnvironmentKVMRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	if c.IsPublic() {
		return resourceEnvironmentKVMReadPublic(d, m)
	}
	var diags diag.Diagnostics
	envName, name := client.KVMDecodeId(d.Id())

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
	retVal := &client.KVM{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)
	d.Set("encrypted", retVal.Encrypted)
	entries := map[string]string{}
	for _, e := range retVal.Entries {
		entries[e.Name] = e.Value
	}
	if retVal.Encrypted {
		d.Set("sensitive_entry", entries)
	} else {
		d.Set("entry", entries)
	}
	return diags
}

func resourceEnvironmentKVMReadPublic(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.KVMDecodeId(d.Id())
	c := m.(*client.Client)

	requestPath := fmt.Sprintf(client.EnvironmentKVMPath, c.Organization, envName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}

	var retVal []string
	err = json.NewDecoder(body).Decode(&retVal)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the specified KVM in the list, no get by ID method exists on public:
	// https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.keyvaluemaps/list
	found := false
	for _, s := range retVal {
		if s == name {
			found = true
			d.Set("environment_name", envName)
			d.Set("name", name)
			d.Set("encrypted", true) // All Apigee cloud offerings are encrypted now

			// Retrieve individual KV's
			requestPath = fmt.Sprintf(client.EnvironmentKVMPathEntries, c.Organization, envName, name)
			body, err = c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
			if err != nil {
				return diag.FromErr(err)
			}

			var res client.KVMEntries
			err = json.NewDecoder(body).Decode(&res)
			if err != nil {
				return diag.FromErr(err)
			}

			entries := map[string]string{}
			for _, e := range res.KeyValueEntries {
				entries[e.Name] = e.Value
			}
			d.Set("sensitive_entry", entries)
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return diags
}

func resourceEnvironmentKVMUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.KVMDecodeId(d.Id())
	c := m.(*client.Client)
	//All other properties besides entries are ForceNew so just handle entries here
	//Check for removal of entries
	encrypted := d.Get("encrypted").(bool)
	var o interface{}
	var n interface{}
	if encrypted {
		o, n = d.GetChange("sensitive_entry")
	} else {
		o, n = d.GetChange("entry")
	}
	oldE := o.(map[string]interface{})
	newE := n.(map[string]interface{})
	for oldKey := range oldE {
		_, newHasKey := newE[oldKey]
		if newHasKey {
			continue
		}
		//Delete entry
		requestPath := fmt.Sprintf(client.EnvironmentKVMPathEntriesGet, c.Organization, envName, name, oldKey)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	//Public Apigee requires entries to be added/changed individually
	if c.IsPublic() {
		//Check for addition/modification of entries
		for newKey := range newE {
			_, oldHasKey := oldE[newKey]
			newOrModValue := newE[newKey].(string)
			//Skip if change with same value
			if oldHasKey {
				oldValue := oldE[newKey].(string)
				if oldValue == newOrModValue {
					continue
				}
			}
			buf := bytes.Buffer{}
			newOrModEntry := client.Attribute{
				Name:  newKey,
				Value: newOrModValue,
			}
			err := json.NewEncoder(&buf).Encode(newOrModEntry)
			if err != nil {
				return diag.FromErr(err)
			}
			requestPath := ""
			if oldHasKey {
				//Change entry
				requestPath = fmt.Sprintf(client.EnvironmentKVMPathEntriesGet, c.Organization, envName, name, newKey)
			} else {
				//Add entry
				requestPath = fmt.Sprintf(client.EnvironmentKVMPathEntries, c.Organization, envName, name)
			}
			_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		buf := bytes.Buffer{}
		upEnvironmentKVM := client.KVM{
			EnvironmentName: envName,
			Name:            name,
		}
		fillEnvironmentKVM(&upEnvironmentKVM, c.IsPublic(), d)
		err := json.NewEncoder(&buf).Encode(upEnvironmentKVM)
		if err != nil {
			return diag.FromErr(err)
		}
		requestPath := fmt.Sprintf(client.EnvironmentKVMPathGet, c.Organization, envName, name)
		_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceEnvironmentKVMDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.KVMDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.EnvironmentKVMPathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
