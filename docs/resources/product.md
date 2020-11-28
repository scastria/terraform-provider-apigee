# Resource: apigee_product
Represents an API product
## Example usage
```hcl
resource "apigee_product" "example" {
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
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of product.
* `display_name` - **(Required, String)** The display name of product.
* `auto_approval_type` - **(Required, Boolean)** Flag that specifies how API keys are approved to access the APIs defined by the API product.
* `description` - **(Optional, String)** The description of product.
* `quota` - **(Optional, Integer)** Number of request messages permitted per app by this API product for the specified `quota_interval` and `quota_time_unit`.
* `quota_interval` - **(Optional, Integer)** Time interval over which the number of request messages is calculated.
* `quota_time_unit` - **(Optional, String)** Time unit defined for the `quota_interval`.  Allowed values: `minute`, `hour`, `day`, `month`. 
* `api_resources` - **(Optional, List of String)** API resources to be bundled in the API product. You can select a specific path, or you can select all subpaths with a wildcard (`/**` and `/*`). 
* `environments` - **(Optional, List of String)** Environment names to which the API product is bound.
* `proxies` - **(Optional, List of String)** API proxy names to which this API product is bound.
* `scopes` - **(Optional, List of String)** OAuth scopes that are validated at runtime.
* `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the product. Use this property to specify the `access` level of the API product as either `public`, `private`, or `internal`.
## Attribute Reference
* `id` - Same as `email`
## Import
Products can be imported using a proper value of `id` as described above
