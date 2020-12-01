---
subcategory: "Publish"
---
# Resource: apigee_company_developer
Represents a developer belonging to a company
## Example usage
```hcl
resource "apigee_company" "MyCompany" {
  name = "MyCompany"
  display_name = "My Company"
  attributes = {
    first = "firstValue"
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
resource "apigee_company_developer" "example" {
  company_name = apigee_company.MyCompany.name
  developer_email = apigee_developer.MyDeveloper.email
  role_name = "myrole"
}
```
## Argument Reference
* `company_name` - **(Required, ForceNew, String)** The name of company
* `developer_email` - **(Required, ForceNew, String)** The email of developer
* `role_name` - **(Optional, String)** The name of a role. This is NOT an `apigee_role` name, but some user defined concept.
## Attribute Reference
* `id` - Same as `company_name`:`developer_email`
## Import
Company developers can be imported using a proper value of `id` as described above
