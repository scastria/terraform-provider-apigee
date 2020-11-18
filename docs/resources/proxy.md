# Resource: apigee_proxy
Represents a proxy's latest revision
## Example usage
```hcl
resource "apigee_proxy" "example" {
  name = "ShawnTest"
  bundle = "proxies/ShawnTest/ShawnTest.zip"
  bundle_hash = filebase64sha256("proxies/ShawnTest/ShawnTest.zip")
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of the proxy.
* `bundle` - **(Required, String)** The filename of the bundle zip.
* `bundle_hash` - **(Required, String)** The hash of the bundle zip used to detect changes of the contents of the zip.
## Attribute Reference
* `id` - Same as `name`
* `revision` - The last revision imported
## Import
Proxies can be imported using a proper value of `id` as described above
