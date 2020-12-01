# Resource: apigee_company_app
Represents an app belonging to a company
## Example usage
```hcl
resource "apigee_company" "MyCompany" {
  name = "MyCompany"
  display_name = "My Company"
  attributes = {
    hello = "goodbye"
  }
}
resource "apigee_company_app" "example" {
  company_name = apigee_company.MyCompany.name
  name = "MyApp"
  callback_url = "hello.com"
  attributes = {
    hello = "goodbye"
  }
}
```
## Argument Reference
* `company_name` - **(Required, ForceNew, String)** The name of a company.
* `name` - **(Required, ForceNew, String)** The name of the app.
* `callback_url` - **(Optional, String)** The callback URL of the app used in OAuth 2.0 authorization code flows.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the app.
## Attribute Reference
* `id` - Same as `company_name`:`name`
## Import
Company apps can be imported using a proper value of `id` as described above
