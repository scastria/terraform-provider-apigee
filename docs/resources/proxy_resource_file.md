---
subcategory: "Develop"
---
# Resource: apigee_proxy_resource_file
Represents a resource file in a proxy
## Example usage
```hcl
resource "apigee_proxy" "MyProxy" {
  name = "MyProxy"
  bundle = "proxies/MyProxy/MyProxy.zip"
  bundle_hash = filebase64sha256("proxies/MyProxy/MyProxy.zip")
}
resource "apigee_proxy_resource_file" "example" {
  proxy_name = apigee_proxy.MyProxy.name
  revision = apigee_proxy.MyProxy.revision
  type = "js"
  name = "test.js"
  file = "resourceFiles/test.js"
  file_hash = filebase64sha256("resourceFiles/test.js")
}
```
## Argument Reference
* `proxy_name` - **(Required, ForceNew, String)** The name of a proxy.
* `revision` - **(Required, ForceNew, Integer)** The revision of a proxy.
* `type` - **(Required, ForceNew, String)** The type of the resource.
* `name` - **(Required, ForceNew, String)** The name of the resource.
* `file` - **(Required, String)** The filename of the resource.
* `file_hash` - **(Required, String)** The hash of the file used to detect changes of the contents of the file.
## Attribute Reference
* `id` - Same as `proxy_name`:`revision`:`type`:`name`
## Import
Proxy resource files can be imported using a proper value of `id` as described above
