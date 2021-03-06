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

func resourceSharedFlowDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFlowDeploymentCreate,
		ReadContext:   resourceSharedFlowDeploymentRead,
		UpdateContext: resourceSharedFlowDeploymentUpdate,
		DeleteContext: resourceSharedFlowDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"shared_flow_name": {
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
				DiffSuppressFunc: resourceSharedFlowDelayDiff,
			},
		},
	}
}

func resourceSharedFlowDelayDiff(k string, old string, n string, d *schema.ResourceData) bool {
	//Suppress all diffs
	return true
}

func resourceSharedFlowDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newSharedFlowDeployment := client.SharedFlowDeployment{
		EnvironmentName: d.Get("environment_name").(string),
		SharedFlowName:  d.Get("shared_flow_name").(string),
	}
	revision := d.Get("revision").(int)
	requestPath := fmt.Sprintf(client.SharedFlowDeploymentRevisionPath, c.Organization, newSharedFlowDeployment.EnvironmentName, newSharedFlowDeployment.SharedFlowName, revision)
	_, err := c.HttpRequest(http.MethodPost, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(newSharedFlowDeployment.SharedFlowDeploymentEncodeId())
	return diags
}

func resourceSharedFlowDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, sharedFlowName := client.SharedFlowDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.SharedFlowDeploymentPath, c.Organization, envName, sharedFlowName)
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
		retVal = &client.GoogleSharedFlowDeployment{}
	} else {
		retVal = &client.SharedFlowDeployment{}
	}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("shared_flow_name", sharedFlowName)
	lastRevision := ""
	//Retrieve the latest revision deployed as THE revision, assumes array is sorted
	if c.IsGoogle() {
		googleRetVal := retVal.(*client.GoogleSharedFlowDeployment)
		lastRevision = googleRetVal.Deployments[len(googleRetVal.Deployments)-1].Revision
	} else {
		oldRetVal := retVal.(*client.SharedFlowDeployment)
		lastRevision = oldRetVal.Revisions[len(oldRetVal.Revisions)-1].Name
	}
	revision, _ := strconv.Atoi(lastRevision)
	d.Set("revision", revision)
	return diags
}

func resourceSharedFlowDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, sharedFlowName := client.SharedFlowDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	revision := d.Get("revision").(int)
	delay := d.Get("delay").(int)
	requestPath := fmt.Sprintf(client.SharedFlowDeploymentRevisionPath, c.Organization, envName, sharedFlowName, revision)
	requestForm := url.Values{
		"override": []string{strconv.FormatBool(true)},
	}
	if !c.IsGoogle() {
		requestForm["delay"] = []string{strconv.Itoa(delay)}
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

func resourceSharedFlowDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	envName, sharedFlowName := client.SharedFlowDeploymentDecodeId(d.Id())
	c := m.(*client.Client)
	//Get all deployments of this shared flow to this environment
	requestPath := fmt.Sprintf(client.SharedFlowDeploymentPath, c.Organization, envName, sharedFlowName)
	body, err := c.HttpRequest(http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	var envDeployments interface{}
	if c.IsGoogle() {
		envDeployments = &client.GoogleSharedFlowDeployment{}
	} else {
		envDeployments = &client.SharedFlowDeployment{}
	}
	err = json.NewDecoder(body).Decode(envDeployments)
	if err != nil {
		return diag.FromErr(err)
	}
	var deployedRevisions []int
	if c.IsGoogle() {
		googleEnvDeployments := envDeployments.(*client.GoogleSharedFlowDeployment)
		for _, dr := range googleEnvDeployments.Deployments {
			revision, _ := strconv.Atoi(dr.Revision)
			deployedRevisions = append(deployedRevisions, revision)
		}
	} else {
		oldEnvDeployments := envDeployments.(*client.SharedFlowDeployment)
		for _, rev := range oldEnvDeployments.Revisions {
			revision, _ := strconv.Atoi(rev.Name)
			deployedRevisions = append(deployedRevisions, revision)
		}
	}
	//Delete each deployment
	for _, revision := range deployedRevisions {
		requestPath := fmt.Sprintf(client.SharedFlowDeploymentRevisionPath, c.Organization, envName, sharedFlowName, revision)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
