# Resource: apigee_role_permission
Represents a permission assigned to a role
## Example usage
```hcl
resource "apigee_role_permission" "example" {
  role_name = "Readers"
  path = "/applications"
  permissions = [
    "get"
  ]
}
```
## Argument Reference
* `role_name` - **(Required, ForceNew, String)** The name of role
* `path` - **(Required, ForceNew, String)** The URI of a permission. See [Apigee Permissions Reference](https://docs.apigee.com/api-platform/system-administration/permissions) for details
* `permissions` - **(Required, ForceNew, List of String)** Any combination of `get`, `put`, and `delete`
## Attribute Reference
* `id` - Same as `role_name`:`path`
## Import
Role permissions can be imported using a proper value of `id` as described above
