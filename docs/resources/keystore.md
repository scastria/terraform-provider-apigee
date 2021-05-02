---
subcategory: "Admin"
---
# Resource: apigee_keystore
Represents a keystore or truststore in an environment
## Example usage
```hcl
resource "apigee_keystore" "example" {
  environment_name = "dev"
  name = "keystoreName"
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the keystore or truststore
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Keystores can be imported using a proper value of `id` as described above
