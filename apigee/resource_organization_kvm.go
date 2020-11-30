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

func resourceOrganizationKVM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationKVMCreate,
		ReadContext:   resourceOrganizationKVMRead,
		UpdateContext: resourceOrganizationKVMUpdate,
		DeleteContext: resourceOrganizationKVMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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
		CustomizeDiff: resourceOrganizationKVMCustomDiff,
	}
}

func resourceOrganizationKVMCustomDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
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

func resourceOrganizationKVMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newOrganizationKVM := client.KVM{
		Name: d.Get("name").(string),
	}
	fillOrganizationKVM(&newOrganizationKVM, d)
	err := json.NewEncoder(&buf).Encode(newOrganizationKVM)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.OrganizationKVMPath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newOrganizationKVM.Name)
	return diags
}

func fillOrganizationKVM(c *client.KVM, d *schema.ResourceData) {
	encrypted, ok := d.GetOk("encrypted")
	if ok {
		c.Encrypted = encrypted.(bool)
	}
	e, ok := d.GetOk("entry")
	if ok {
		entries := e.(map[string]interface{})
		for name, value := range entries {
			c.Entries = append(c.Entries, client.Attribute{
				Name:  name,
				Value: value.(string),
			})
		}
	}
}

func resourceOrganizationKVMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.OrganizationKVMPathGet, c.Organization, d.Id())
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
	d.Set("name", d.Id())
	entries := map[string]string{}
	for _, e := range retVal.Entries {
		entries[e.Name] = e.Value
	}
	d.Set("entry", entries)
	return diags
}

func resourceOrganizationKVMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	//Check for removal of entries
	if d.HasChange("entry") {
		o, n := d.GetChange("entry")
		oldE := o.(map[string]interface{})
		newE := n.(map[string]interface{})
		for oldKey, _ := range oldE {
			_, newHasKey := newE[oldKey]
			if newHasKey {
				continue
			}
			//Delete entry
			requestPath := fmt.Sprintf(client.OrganizationKVMPathGetEntry, c.Organization, d.Id(), oldKey)
			_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	buf := bytes.Buffer{}
	upOrganizationKVM := client.KVM{
		Name: d.Id(),
	}
	fillOrganizationKVM(&upOrganizationKVM, d)
	err := json.NewEncoder(&buf).Encode(upOrganizationKVM)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.OrganizationKVMPathGet, c.Organization, d.Id())
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceOrganizationKVMDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.OrganizationKVMPathGet, c.Organization, d.Id())
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
