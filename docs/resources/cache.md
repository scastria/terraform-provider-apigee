# Resource: apigee_cache
Represents a cache in an environment
## Example usage
```hcl
resource "apigee_cache" "example" {
  environment_name = "dev"
  name = "Tokens"
  description = "OIDC access tokens"
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the cache
* `description` - **(Optional, String)** The description of the cache
* `expiry_timeout_in_sec` - **(Optional, String)** The description of the cache
## Attribute Reference
* `id` - Same as `role_name`:`path`
## Import
Role permissions can be imported using a proper value of `id` as described above
