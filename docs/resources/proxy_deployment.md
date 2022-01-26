---
subcategory: "Develop"
---
# Resource: apigee_proxy_deployment
Represents a proxy's deployment to an environment
## Example usage
```hcl
resource "apigee_proxy" "MyProxy" {
  name = "MyProxy"
  bundle = "proxies/MyProxy/MyProxy.zip"
  bundle_hash = filebase64sha256("proxies/MyProxy/MyProxy.zip")
}
resource "apigee_proxy_deployment" "example" {
  proxy_name = apigee_proxy.MyProxy.name
  environment_name = "dev"
  revision = apigee_proxy.MyProxy.revision
}
```
## Argument Reference
* `proxy_name` - **(Required, ForceNew, String)** The name of the proxy to be deployed.
* `environment_name` - **(Required, ForceNew, String)** The environment to deploy the proxy to.
* `revision` - **(Required, Integer)** The revision of the proxy to deploy.  On create, it will assume the proxy has not been deployed in the given environment yet.  On update, it will override any current deployment to the given environment.
* `delay` - **(Optional, Integer)** Time interval, in seconds, to wait before undeploying the currently deployed revision.  Default: 0. Ignored for calculating diffs.
* `service_account` - **(Optional, String)** For Google Cloud Apigee version, specify the service account associated with the deployment. See the [Google documentation](https://cloud.google.com/apigee/docs/api-platform/security/google-auth/overview#about-service-account-permissions) for permissions required by the deploying user.
## Attribute Reference
* `id` - Same as `environment_name`:`proxy_name`
## Import
Proxy deployments can be imported using a proper value of `id` as described above
