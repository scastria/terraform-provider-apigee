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
resource "apigee_proxy_kvm" "encryptedExample" {
  proxy_name = apigee_proxy.MyProxy.name
  name = "LookupValues"
  encrypted = true
  sensitive_entry = {
    first = "firstValue"
    second = "secondValue"
  }
}
```
## Argument Reference
* `proxy_name` - **(Required, ForceNew, String)** The name of a proxy
* `name` - **(Required, ForceNew, String)** The name of the kvm
* `encrypted` - **(Optional, ForceNew, Boolean)** Determine whether to encrypt the values within the kvm.  Due to Apigee API, encrypted values can NOT be read back, therefore, a change will always be detected even when there may not be one.  You can use `lifecycle` and `ignore_changes` to avoid this issue.
* `entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm when `encrypted` is `false`.  Values will NOT be hidden from logs.
* `sensitive_entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm when `encrypted` is `true`.  Values WILL be hidden from logs.
## Attribute Reference
* `id` - Same as `proxy_name`:`name`
## Import
Proxy KVMs can be imported using a proper value of `id` as described above
