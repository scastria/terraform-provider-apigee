---
subcategory: "Publish"
---
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
  operation {
    api_source = "proxy_name"
    path       = "/v1"
    methods    = ["GET"]

    quota           = 5000
    quota_interval  = 1
    quota_time_unit = "month"

    attributes = {
      message-weight = "1"
    }
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
* `operations` - **(Optional, Block of Operations)** A list of Operations supported by this product.
  * `api_source` - **(Required, String)** The name of the Apigee Proxy see [proxy](proxy.md)
  * `path` - **(Required, String)** The path this product can request e.g. /v1/**
  * `methods` - **(Required, List String)** Supported HTTP methods e.g. ["GET", "POST"]
  * `quota` - **(Optional, Integer)** Number of request messages permitted per app by this API product for the specified `quota_interval` and `quota_time_unit`.
  * `quota_interval` - **(Optional, Integer)** Time interval over which the number of request messages is calculated.
  * `quota_time_unit` - **(Optional, String)** Time unit defined for the `quota_interval`.  Allowed values: `minute`, `hour`, `day`, `month`.
  * `attributes` - **(Optional, Map of String to String)** Keys and values to be stored as custom attributes of the operation.
* `operation_config_type` - **(Optional, String)** The Operation config type for the product, can either be `proxy` or `remoteservice`.

## Attribute Reference
* `id` - Same as `name`

## Import
Products can be imported using a proper value of `id` as described above
