# Resource: apigee_developer_app
Represents an app belonging to a developer
## Example usage
```hcl
resource "apigee_developer" "MyDeveloper" {
  email = "ahamilton@example.com"
  first_name = "Alex"
  last_name = "Hamilton"
  user_name = "ahamilton@example.com"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_developer_app" "example" {
  developer_email = apigee_developer.MyDeveloper.email
  name = "MyApp"
  callback_url = "hello.com"
  attributes = {
    hello = "goodbye"
  }
}
```
## Argument Reference
* `developer_email` - **(Required, ForceNew, String)** The email address of a developer.
* `name` - **(Required, ForceNew, String)** The name of the app.
* `callback_url` - **(Optional, String)** The callback URL of the app used in OAuth 2.0 authorization code flows.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the app.
## Attribute Reference
* `id` - Same as `developer_email`:`name`
## Import
Developer apps can be imported using a proper value of `id` as described above
