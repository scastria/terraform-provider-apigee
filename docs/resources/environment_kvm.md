---
subcategory: "Admin"
---
# Resource: apigee_environment_kvm
Represents a kvm in an environment
## Example usage
```hcl
resource "apigee_environment_kvm" "example" {
  environment_name = "dev"
  name = "LookupValues"
  entry = {
    first = "firstValue"
    second = "secondValue"
  }
}
resource "apigee_environment_kvm" "encryptedExample" {
  environment_name = "dev"
  name = "LookupValues"
  encrypted = true
  sensitive_entry = {
    first = "firstValue"
    second = "secondValue"
  }
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the kvm
* `encrypted` - **(Optional, ForceNew, Boolean)** Determine whether to encrypt the values within the kvm.  Due to Apigee API, encrypted values can NOT be read back, therefore, a change will always be detected even when there may not be one.  You can use `lifecycle` and `ignore_changes` to avoid this issue. 
* `entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm when `encrypted` is `false`.  Values will NOT be hidden from logs.
* `sensitive_entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm when `encrypted` is `true`.  Values WILL be hidden from logs.
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Environment KVMs can be imported using a proper value of `id` as described above
