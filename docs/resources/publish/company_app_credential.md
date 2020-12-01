# Resource: apigee_company_app_credential
Represents a credential belonging to a company app
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
resource "apigee_company" "MyCompany" {
  name = "MyCompany"
  display_name = "My Company"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_company_app" "MyApp" {
  company_name = apigee_company.MyCompany.name
  name = "MyApp"
  callback_url = "hello.com"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_company_app_credential" "example" {
  company_name = apigee_company.MyCompany.name
  company_app_name = apigee_company_app.MyApp.name
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
* `company_name` - **(Required, ForceNew, String)** The name of a company.
* `company_app_name` - **(Required, ForceNew, String)** The name of a company app.
* `consumer_key` - **(Required, ForceNew, String)** The key of credential.
* `consumer_secret` - **(Required, ForceNew, String)** The secret of credential.
* `api_products` - **(Optional, List of String)** The API products to associate this credential with.
* `scopes` - **(Optional, List of String)** The scopes to allow this credential to be used with.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the credential.
## Attribute Reference
* `id` - Same as `company_name`:`company_app_name`:`consumer_key`
## Import
Company app credentials can be imported using a proper value of `id` as described above
