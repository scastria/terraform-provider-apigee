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
	"net/http"
)

func resourceKeystore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeystoreCreate,
		ReadContext:   resourceKeystoreRead,
		DeleteContext: resourceKeystoreDelete,
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
		},
	}
}

func resourceKeystoreCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newKeystore := client.Keystore{
		EnvironmentName: d.Get("environment_name").(string),
		Name:            d.Get("name").(string),
	}
	err := json.NewEncoder(&buf).Encode(newKeystore)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.KeystorePath, c.Organization, newKeystore.EnvironmentName)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newKeystore.KeystoreEncodeId())
	return diags
}

func resourceKeystoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.KeystoreDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.KeystorePathGet, c.Organization, envName, name)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Keystore{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("name", name)

	return diags
}

func resourceKeystoreDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, name := client.KeystoreDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.KeystorePathGet, c.Organization, envName, name)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
