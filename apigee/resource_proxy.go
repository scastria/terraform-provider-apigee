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
	"net/url"
	"strconv"
)

func resourceProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyCreate,
		ReadContext:   resourceProxyRead,
		UpdateContext: resourceProxyUpdate,
		DeleteContext: resourceProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bundle": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bundle_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		CustomizeDiff: resourceProxyCustomDiff,
	}
}

func resourceProxyCustomDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	//Mark the revision as changing if bundle changes
	if diff.HasChange("bundle") {
		diff.SetNewComputed("revision")
	}
	if diff.HasChange("bundle_hash") {
		diff.SetNewComputed("revision")
	}
	return nil
}

func importProxyRevision(c *client.Client, name string, bundle string) (*client.ProxyRevision, error) {
	//Turn filename into multi part buffer
	mp, buf, err := client.GetMultiPartBuffer(map[string]client.FormData{
		"bundle": client.FormData{Filename: bundle},
	})
	if err != nil {
		return nil, err
	}
	requestPath := fmt.Sprintf(client.ProxyPath, c.Organization)
	requestHeaders := http.Header{
		headers.ContentType: []string{mp.FormDataContentType()},
	}
	requestQuery := url.Values{
		"action": []string{"import"},
		"name":   []string{name},
	}
	body, err := c.HttpRequest(http.MethodPost, requestPath, requestQuery, requestHeaders, buf)
	if err != nil {
		return nil, err
	}
	retVal := &client.ProxyRevision{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

func resourceProxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Get("name").(string)
	bundle := d.Get("bundle").(string)
	retVal, err := importProxyRevision(c, name, bundle)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(name)
	revision, _ := strconv.Atoi(retVal.Revision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyPathGet, c.Organization, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Proxy{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("name", d.Id())
	//Empirically the revisions array is sorted *alphabetically* when we get
	//it which means that, for example, an API with 10 revisions comes back
	//with a revision list of 1, 10, 2, 3, 4, 5, 6, 7, 8, 9. As such we need
	//to iteratate over the entire array to determine the latest revision.
	//Any errors which arise parsing revision numbers are returned as a
	//diagnostic.
	revision := 0
	for _, revisionStr := range retVal.Revisions {
		rn, err := strconv.Atoi(revisionStr)
		if err != nil {
			return diag.FromErr(err)
		}
		if rn > revision {
			revision = rn
		}
	}
	if revision == 0 {
		return diag.Errorf("proxy has no latest revision")
	}
	d.Set("revision", revision)
	return diags
}

func resourceProxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	bundle := d.Get("bundle").(string)
	retVal, err := importProxyRevision(c, d.Id(), bundle)
	if err != nil {
		return diag.FromErr(err)
	}
	revision, _ := strconv.Atoi(retVal.Revision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	//Assume that if this proxy is currently deployed to ANY environment, that it was deployed in a different
	//TF configuration.  Therefore, that other TF configuration will handle the delete of the proxy so just
	//report deletion to TF even though it was not deleted from Apigee.  This is a design decision which
	//could be solved other ways, but would be much more complicated for the user.
	//Get all deployments of this proxy to ANY environment
	requestPath := fmt.Sprintf(client.ProxyDeploymentPath, c.Organization, d.Id())
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	var deployments interface{}
	if c.IsGoogle() {
		deployments = &client.GoogleProxyEnvironmentDeployment{}
	} else {
		deployments = &client.ProxyDeployments{}
	}
	err = json.NewDecoder(body).Decode(deployments)
	if err != nil {
		return diag.FromErr(err)
	}
	var numDeployments int
	if c.IsGoogle() {
		googleDeployments := deployments.(*client.GoogleProxyEnvironmentDeployment)
		numDeployments = len(googleDeployments.Deployments)
	} else {
		oldDeployments := deployments.(*client.ProxyDeployments)
		numDeployments = len(oldDeployments.Environments)
	}
	if numDeployments == 0 {
		requestPath = fmt.Sprintf(client.ProxyPathGet, c.Organization, d.Id())
		_, err = c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
