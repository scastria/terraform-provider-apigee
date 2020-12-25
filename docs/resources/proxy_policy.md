---
subcategory: "Develop"
---
# Resource: apigee_proxy_policy
Represents a policy in a proxy
## Example usage
```hcl
resource "apigee_proxy" "MyProxy" {
  name = "MyProxy"
  bundle = "proxies/MyProxy/MyProxy.zip"
  bundle_hash = filebase64sha256("proxies/MyProxy/MyProxy.zip")
}
resource "apigee_proxy_policy" "example" {
  proxy_name = apigee_proxy.MyProxy.name
  revision = apigee_proxy.MyProxy.revision
  name = "VerifyAccessToken"
  file = "policies/test.xml"
  file_hash = filebase64sha256("policies/test.xml")
}
```
## Argument Reference
* `proxy_name` - **(Required, ForceNew, String)** The name of a proxy.
* `revision` - **(Required, ForceNew, Integer)** The revision of a proxy.
* `name` - **(Required, ForceNew, String)** The name of the policy.
* `file` - **(Required, ForceNew, String)** The filename of the policy.
* `file_hash` - **(Required, ForceNew, String)** The hash of the file used to detect changes of the contents of the file.
## Attribute Reference
* `id` - Same as `proxy_name`:`revision`:`name`
## Import
Proxy policies can be imported using a proper value of `id` as described above
