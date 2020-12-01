---
subcategory: "Develop"
---
# Resource: apigee_proxy_kvm
Represents a kvm in a proxy
## Example usage
```hcl
resource "apigee_proxy" "MyProxy" {
  name = "MyProxy"
  bundle = "proxies/MyProxy/MyProxy.zip"
  bundle_hash = filebase64sha256("proxies/MyProxy/MyProxy.zip")
}
resource "apigee_proxy_kvm" "example" {
  proxy_name = apigee_proxy.MyProxy.name
  name = "LookupValues"
  entry = {
    first = "firstValue"
    second = "secondValue"
  }
}
```
## Argument Reference
* `proxy_name` - **(Required, ForceNew, String)** The name of a proxy
* `name` - **(Required, ForceNew, String)** The name of the cache
* `encrypted` - **(Optional, Boolean)** Determine whether to encrypt the values within the kvm.  Changing this value from `true` to `false` will cause ForceNew since Apigee will not decrypt values. 
* `entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm.
## Attribute Reference
* `id` - Same as `proxy_name`:`name`
## Import
Proxy KVMs can be imported using a proper value of `id` as described above
