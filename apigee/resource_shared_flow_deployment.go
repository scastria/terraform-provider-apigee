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
	retVal := &client.SharedFlowDeployment{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.Set("environment_name", envName)
	d.Set("shared_flow_name", sharedFlowName)
	//Retrieve the latest revision deployed as THE revision, assumes array is sorted
	lastRevision := retVal.Revisions[len(retVal.Revisions)-1]
	revision, _ := strconv.Atoi(lastRevision.Name)
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
	envDeployments := &client.SharedFlowDeployment{}
	err = json.NewDecoder(body).Decode(envDeployments)
	if err != nil {
		return diag.FromErr(err)
	}
	//Delete each deployment
	for _, rev := range envDeployments.Revisions {
		revision, _ := strconv.Atoi(rev.Name)
		requestPath := fmt.Sprintf(client.SharedFlowDeploymentRevisionPath, c.Organization, envName, sharedFlowName, revision)
		_, err := c.HttpRequest(http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
