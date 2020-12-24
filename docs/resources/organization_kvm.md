---
subcategory: "Admin"
---
# Resource: apigee_organization_kvm
Represents a kvm in an organization
## Example usage
```hcl
resource "apigee_organization_kvm" "example" {
  name = "LookupValues"
  entry = {
    first = "firstValue"
    second = "secondValue"
  }
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of the kvm
* `encrypted` - **(Optional, Boolean)** Determine whether to encrypt the values within the kvm.  Changing this value from `true` to `false` will cause ForceNew since Apigee will not decrypt values. 
* `entry` - **(Optional, Map of String to String)** Keys and values to be stored within the kvm.
## Attribute Reference
* `id` - Same as `name`
## Import
Organization KVMs can be imported using a proper value of `id` as described above
