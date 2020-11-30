# Resource: apigee_developer_app_credential
Represents a credential belonging to a developer app
## Example usage
```hcl
resource "apigee_product" "MyProduct" {
  name = "MyProduct"
  display_name = "MyProduct"
  auto_approval_type = true
  description = "A great product"
  environments = [
    "dev",
    "test",
    "stage"
  ]
  scopes = [
    "openid",
    "profile"
  ]
  attributes = {
    access = "public"
  }
}
resource "apigee_developer" "MyDeveloper" {
  email = "ahamilton@example.com"
  first_name = "Alex"
  last_name = "Hamilton"
  user_name = "ahamilton@example.com"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_developer_app" "MyApp" {
  developer_email = apigee_developer.MyDeveloper.email
  name = "MyApp"
  callback_url = "hello.com"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_developer_app_credential" "example" {
  developer_email = apigee_developer.MyDeveloper.email
  developer_app_name = apigee_developer_app.MyApp.name
  consumer_key = "MyKey"
  consumer_secret = "secret"
  api_products = [
    apigee_product.MyProduct.name
  ]
  scopes = [
    "openid"
  ]
  attributes = {
    hello = "goodbye"
  }
}
```
## Argument Reference
* `developer_email` - **(Required, ForceNew, String)** The email address of a developer. Can be changed to rename developer.
* `developer_app_name` - **(Required, ForceNew, String)** The name of a developer app.
* `consumer_key` - **(Required, ForceNew, String)** The key of credential.
* `consumer_secret` - **(Required, ForceNew, String)** The secret of credential.
* `api_products` - **(Optional, List of String)** The API products to associate this credential with.
* `scopes` - **(Optional, List of String)** The scopes to allow this credential to be used with.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the credential.
## Attribute Reference
* `id` - Same as `developer_email`:`developer_app_name`:`consumer_key`
## Import
Developer app credentials can be imported using a proper value of `id` as described above
