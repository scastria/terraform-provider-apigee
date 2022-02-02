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
	"strings"
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
			"service_account": {
				Type:     schema.TypeString,
				Optional: true,
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
	v := url.Values{}
	newProxyDeployment := client.ProxyEnvironmentDeployment{
		EnvironmentName: d.Get("environment_name").(string),
		ProxyName:       d.Get("proxy_name").(string),
	}
	if d.Get("service_account").(string) != "" {
		if !c.IsGoogle() {
			return diag.Errorf("service_account cannot be set for non-Google Cloud Apigee versions")
		}
		newProxyDeployment.ServiceAccount = d.Get("service_account").(string)
		v.Set("serviceAccount", newProxyDeployment.ServiceAccount)
	}
	revision := d.Get("revision").(int)
	requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentRevisionPath, c.Organization, newProxyDeployment.EnvironmentName, newProxyDeployment.ProxyName, revision)
	_, err := c.HttpRequest(http.MethodPost, requestPath, v, nil, &bytes.Buffer{})
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
	requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentPath, c.Organization, envName, proxyName)
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
		retVal = &client.GoogleProxyEnvironmentDeployment{}
	} else {
		retVal = &client.ProxyEnvironmentDeployment{}
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
		googleRetVal := retVal.(*client.GoogleProxyEnvironmentDeployment)
		lastRevision = googleRetVal.Deployments[len(googleRetVal.Deployments)-1].Revision
	} else {
		oldRetVal := retVal.(*client.ProxyEnvironmentDeployment)
		lastRevision = oldRetVal.Revisions[len(oldRetVal.Revisions)-1].Name
	}
	revision, _ := strconv.Atoi(lastRevision)
	d.Set("revision", revision)
	if c.IsGoogle() {
		googleRetVal := retVal.(*client.GoogleProxyEnvironmentDeployment)
		serviceAccount := googleRetVal.Deployments[len(googleRetVal.Deployments)-1].ServiceAccount
		//When reading the service account, it is prefixed by "projects/-/serviceAccounts/"
		d.Set("service_account", strings.TrimPrefix(serviceAccount, "projects/-/serviceAccounts/"))
	} else {
		d.Set("service_account", "")
	}
	return diags
}

func resourceProxyDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, proxyName := client.ProxyDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	//Grab revision of existing deployment in case I need to manually undeploy later if the basepath has changed
	prevRevision := 0
	if d.HasChange("revision") {
		o, _ := d.GetChange("revision")
		prevRevision = o.(int)
	}
	revision := d.Get("revision").(int)
	delay := d.Get("delay").(int)
	requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentRevisionPath, c.Organization, envName, proxyName, revision)
	requestForm := url.Values{
		"override": []string{strconv.FormatBool(true)},
	}
	if !c.IsGoogle() {
		requestForm["delay"] = []string{strconv.Itoa(delay)}
	}
	if d.Get("service_account").(string) != "" {
		if !c.IsGoogle() {
			return diag.Errorf("service_account cannot be set for non-Google Cloud Apigee versions")
		}
		requestForm["serviceAccount"] = []string{d.Get("service_account").(string)}
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.FormEncoded},
	}
	_, err := c.HttpRequest(http.MethodPost, requestPath, nil, requestHeaders, bytes.NewBufferString(requestForm.Encode()))
	if err != nil {
		return diag.FromErr(err)
	}
	//Force undeployment of previous revision, ignore errors if undeployed by override = true above
	if prevRevision != 0 {
		requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentRevisionPath, c.Organization, envName, proxyName, prevRevision)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			//Ignore errors
		}
	}
	return diags
}

func resourceProxyDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, proxyName := client.ProxyDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	//Get all deployments of this proxy to this environment
	requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentPath, c.Organization, envName, proxyName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	var envDeployments interface{}
	if c.IsGoogle() {
		envDeployments = &client.GoogleProxyEnvironmentDeployment{}
	} else {
		envDeployments = &client.ProxyEnvironmentDeployment{}
	}
	err = json.NewDecoder(body).Decode(envDeployments)
	if err != nil {
		return diag.FromErr(err)
	}
	var deployedRevisions []int
	if c.IsGoogle() {
		googleEnvDeployments := envDeployments.(*client.GoogleProxyEnvironmentDeployment)
		for _, dr := range googleEnvDeployments.Deployments {
			revision, _ := strconv.Atoi(dr.Revision)
			deployedRevisions = append(deployedRevisions, revision)
		}
	} else {
		oldEnvDeployments := envDeployments.(*client.ProxyEnvironmentDeployment)
		for _, rev := range oldEnvDeployments.Revisions {
			revision, _ := strconv.Atoi(rev.Name)
			deployedRevisions = append(deployedRevisions, revision)
		}
	}
	//Delete each deployment
	for _, revision := range deployedRevisions {
		requestPath := fmt.Sprintf(client.ProxyEnvironmentDeploymentRevisionPath, c.Organization, envName, proxyName, revision)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
