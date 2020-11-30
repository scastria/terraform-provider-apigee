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
	"regexp"
)

func resourceDeveloperAppCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeveloperAppCredentialCreate,
		ReadContext:   resourceDeveloperAppCredentialRead,
		UpdateContext: resourceDeveloperAppCredentialUpdate,
		DeleteContext: resourceDeveloperAppCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"developer_email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`), "must be a valid email address"),
			},
			"developer_app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_secret": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"api_products": {
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

func resourceDeveloperAppCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newDeveloperAppCredential := client.DeveloperAppCredentialModify{
		DeveloperEmail:   d.Get("developer_email").(string),
		DeveloperAppName: d.Get("developer_app_name").(string),
		ConsumerKey:      d.Get("consumer_key").(string),
		ConsumerSecret:   d.Get("consumer_secret").(string),
	}
	fillDeveloperAppCredential(&newDeveloperAppCredential, d)
	err := json.NewEncoder(&buf).Encode(newDeveloperAppCredential)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperAppCredentialPathCreate, c.Organization, newDeveloperAppCredential.DeveloperEmail, newDeveloperAppCredential.DeveloperAppName)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Set id before adding products
	d.SetId(newDeveloperAppCredential.DeveloperAppCredentialEncodeId())
	//Add any products or attributes with POST
	requestPath = fmt.Sprintf(client.DeveloperAppCredentialPathGet, c.Organization, newDeveloperAppCredential.DeveloperEmail, newDeveloperAppCredential.DeveloperAppName, newDeveloperAppCredential.ConsumerKey)
	buf = bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(newDeveloperAppCredential)
	if err != nil {
		//Don't clear id since credential was created
		return diag.FromErr(err)
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		//Don't clear id since credential was created
		return diag.FromErr(err)
	}
	//Add any scopes with PUT
	buf = bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(newDeveloperAppCredential)
	if err != nil {
		//Don't clear id since credential was created
		return diag.FromErr(err)
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		//Don't clear id since credential was created
		return diag.FromErr(err)
	}
	return diags
}

func fillDeveloperAppCredential(c *client.DeveloperAppCredentialModify, d *schema.ResourceData) {
	apiProducts, ok := d.GetOk("api_products")
	if ok {
		set := apiProducts.(*schema.Set)
		c.APIProducts = convertSetToArray(set)
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

func resourceDeveloperAppCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, appName, key := client.DeveloperAppCredentialDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperAppCredentialPathGet, c.Organization, developerEmail, appName, key)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.DeveloperAppCredential{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("developer_email", developerEmail)
	d.Set("developer_app_name", appName)
	d.Set("consumer_key", key)
	d.Set("consumer_secret", retVal.ConsumerSecret)
	var apiProducts []string
	for _, prod := range retVal.APIProducts {
		apiProducts = append(apiProducts, prod.APIProduct)
	}
	d.Set("api_products", apiProducts)
	d.Set("scopes", retVal.Scopes)
	atts := map[string]string{}
	for _, e := range retVal.Attributes {
		atts[e.Name] = e.Value
	}
	d.Set("attributes", atts)
	return diags
}

func resourceDeveloperAppCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, appName, key := client.DeveloperAppCredentialDecodeId(d.Id())
	c := m.(*client.Client)
	//Check for removal of products
	if d.HasChange("api_products") {
		o, n := d.GetChange("api_products")
		oldP := convertSetToArray(o.(*schema.Set))
		newP := convertSetToArray(n.(*schema.Set))
		for _, oldProd := range oldP {
			_, newHasProd := find(newP, oldProd)
			if newHasProd {
				continue
			}
			//Delete product
			requestPath := fmt.Sprintf(client.DeveloperAppCredentialPathProduct, c.Organization, developerEmail, appName, key, oldProd)
			_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	buf := bytes.Buffer{}
	upDeveloperAppCredential := client.DeveloperAppCredentialModify{
		DeveloperEmail:   developerEmail,
		DeveloperAppName: appName,
		ConsumerKey:      key,
		ConsumerSecret:   d.Get("consumer_secret").(string),
	}
	fillDeveloperAppCredential(&upDeveloperAppCredential, d)
	//Handle products and attributes with POST
	err := json.NewEncoder(&buf).Encode(upDeveloperAppCredential)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.DeveloperAppCredentialPathGet, c.Organization, developerEmail, appName, key)
	requestHeaders := http.Header{
		headers.ContentType: []string{mime.TypeByExtension(".json")},
	}
	_, err = c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	//Handle scopes with PUT
	buf = bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(upDeveloperAppCredential)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.HttpRequest(http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceDeveloperAppCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	developerEmail, appName, key := client.DeveloperAppCredentialDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.DeveloperAppCredentialPathGet, c.Organization, developerEmail, appName, key)
	_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
