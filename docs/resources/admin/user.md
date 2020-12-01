# Resource: apigee_user
Represents a user
## Example usage
```hcl
resource "apigee_user" "example" {
  email_id = "John.Smith@ihsmarkit.com"
  first_name = "John"
  last_name = "Smith"
  password = "XXXX"
}
```
## Argument Reference
* `email_id` - **(Required, String)** The email address of user. Can be changed to rename user.
* `first_name` - **(Required, String)** The first name of user.
* `last_name` - **(Required, String)** The last name of user.
* `password` - **(Required, String)** The password of user. Cannot be imported as Apigee does not return it in Management API.
## Attribute Reference
* `id` - Same as `email_id`
## Import
Users can be imported using a proper value of `id` as described above
