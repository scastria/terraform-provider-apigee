# Resource: apigee_developer
Represents a developer
## Example usage
```hcl
resource "apigee_developer" "example" {
  email = "ahamilton@example.com"
  first_name = "Alex"
  last_name = "Hamilton"
  user_name = "ahamilton@example.com"
  attributes = {
    hello = "goodbye"
  }
}
```
## Argument Reference
* `email` - **(Required, String)** The email address of developer. Can be changed to rename developer.
* `first_name` - **(Required, String)** The first name of developer.
* `last_name` - **(Required, String)** The last name of developer.
* `user_name` - **(Required, String)** The user name of developer.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the developer.
## Attribute Reference
* `id` - Same as `email`
## Import
Developers can be imported using a proper value of `id` as described above
