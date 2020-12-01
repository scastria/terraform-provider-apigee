---
subcategory: "Admin"
---
# Resource: apigee_role
Represents a role
## Example usage
```hcl
resource "apigee_role" "example" {
  name = "Readers"
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of role
## Attribute Reference
* `id` - Same as `name`
## Import
Roles can be imported using a proper value of `id` as described above
