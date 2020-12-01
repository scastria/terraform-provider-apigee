---
subcategory: "Publish"
---
# Resource: apigee_company
Represents a company
## Example usage
```hcl
resource "apigee_company" "example" {
  name = "MyCompany"
  display_name = "My Company"
  attributes = {
    first = "firstValue"
  }
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of company.
* `display_name` - **(Optional, String)** The display name of company.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the company.
## Attribute Reference
* `id` - Same as `name`
## Import
Companies can be imported using a proper value of `id` as described above
