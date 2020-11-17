# Resource: apigee_user_role
Represents a user assigned to a role
## Example usage
```hcl
resource "apigee_user_role" "example" {
  email_id = "John.Smith@ihsmarkit.com"
  role_name = "Readers"
}
```
## Argument Reference
* `email_id` - **(Required, ForceNew, String)** The email of user
* `role_name` - **(Required, ForceNew, String)** The name of role
## Attribute Reference
* `id` - Same as `email_id`:`role_name`
## Import
User roles can be imported using a proper value of `id` as described above
