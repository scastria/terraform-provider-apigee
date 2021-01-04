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
	"net/url"
	"strconv"
)

func resourceProxyDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProxyDeploymentCreate,
		ReadContext:   resourceProxyDeploymentRead,
		UpdateContext: resourceProxyDeploymentUpdate,
		DeleteContext: resourceProxyDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"proxy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"revision": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"delay": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntAtLeast(0),
				DiffSuppressFunc: resourceProxyDelayDiff,
			},
		},
	}
}

func resourceProxyDelayDiff(k string, old string, n string, d *schema.ResourceData) bool {
	//Suppress all diffs
	return true
}

func resourceProxyDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newProxyDeployment := client.ProxyDeployment{
		EnvironmentName: d.Get("environment_name").(string),
		ProxyName:       d.Get("proxy_name").(string),
	}
	revision := d.Get("revision").(int)
	requestPath := fmt.Sprintf(client.ProxyDeploymentRevisionPath, c.Organization, newProxyDeployment.EnvironmentName, newProxyDeployment.ProxyName, revision)
	_, err := c.HttpRequest(http.MethodPost, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newProxyDeployment.ProxyDeploymentEncodeId())
	return diags
}

func resourceProxyDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, proxyName := client.ProxyDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ProxyDeploymentPath, c.Organization, envName, proxyName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	var retVal interface{}
	if c.IsGoogle() {
		retVal = &client.GoogleProxyDeployment{}
	} else {
		retVal = &client.ProxyDeployment{}
	}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("proxy_name", proxyName)
	lastRevision := ""
	//Retrieve the latest revision deployed as THE revision, assumes array is sorted
	if c.IsGoogle() {
		googleRetVal := retVal.(*client.GoogleProxyDeployment)
		lastRevision = googleRetVal.Deployments[len(googleRetVal.Deployments)-1].Revision
	} else {
		oldRetVal := retVal.(*client.ProxyDeployment)
		lastRevision = oldRetVal.Revisions[len(oldRetVal.Revisions)-1].Name
	}
	revision, _ := strconv.Atoi(lastRevision)
	d.Set("revision", revision)
	return diags
}

func resourceProxyDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, proxyName := client.ProxyDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	revision := d.Get("revision").(int)
	delay := d.Get("delay").(int)
	requestPath := fmt.Sprintf(client.ProxyDeploymentRevisionPath, c.Organization, envName, proxyName, revision)
	requestForm := url.Values{
		"override": []string{strconv.FormatBool(true)},
		"delay":    []string{strconv.Itoa(delay)},
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.FormEncoded},
	}
	_, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, bytes.NewBufferString(requestForm.Encode()))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceProxyDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, proxyName := client.ProxyDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	//Get all deployments of this proxy to this environment
	requestPath := fmt.Sprintf(client.ProxyDeploymentPath, c.Organization, envName, proxyName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	var envDeployments interface{}
	if c.IsGoogle() {
		envDeployments = &client.GoogleProxyDeployment{}
	} else {
		envDeployments = &client.ProxyDeployment{}
	}
	err = json.NewDecoder(body).Decode(envDeployments)
	if err != nil {
		return diag.FromErr(err)
	}
	var deployedRevisions []int
	if c.IsGoogle() {
		googleEnvDeployments := envDeployments.(*client.GoogleProxyDeployment)
		for _, dr := range googleEnvDeployments.Deployments {
			revision, _ := strconv.Atoi(dr.Revision)
			deployedRevisions = append(deployedRevisions, revision)
		}
	} else {
		oldEnvDeployments := envDeployments.(*client.ProxyDeployment)
		for _, rev := range oldEnvDeployments.Revisions {
			revision, _ := strconv.Atoi(rev.Name)
			deployedRevisions = append(deployedRevisions, revision)
		}
	}
	//Delete each deployment
	for _, revision := range deployedRevisions {
		requestPath := fmt.Sprintf(client.ProxyDeploymentRevisionPath, c.Organization, envName, proxyName, revision)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
