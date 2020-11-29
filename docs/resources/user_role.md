# Resource: apigee_user_role
Represents a user assigned to a role
## Example usage
```hcl
resource "apigee_user" "MyUser" {
  email_id = "John.Smith@ihsmarkit.com"
  first_name = "John"
  last_name = "Smith"
  password = "XXXX"
}
resource "apigee_role" "MyRole" {
  name = "Readers"
}
resource "apigee_user_role" "example" {
  email_id = apigee_user.MyUser.email_id
  role_name = apigee_role.MyRole.name
}
```
## Argument Reference
* `email_id` - **(Required, ForceNew, String)** The email of user
* `role_name` - **(Required, ForceNew, String)** The name of role
## Attribute Reference
* `id` - Same as `email_id`:`role_name`
## Import
User roles can be imported using a proper value of `id` as described above
