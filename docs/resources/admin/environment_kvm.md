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
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the cache
* `encrypted` - **(Optional, Boolean)** Determine whether to encrypt the values within the kvm.  Changing this value from `true` to `false` will cause ForceNew since Apigee will not decrypt values. 
* `entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm.
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Environment KVMs can be imported using a proper value of `id` as described above
