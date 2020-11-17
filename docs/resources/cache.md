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
* `expiry_timeout_in_sec` - **(Optional, Integer)** The default timeout in seconds of entries within the cache.  Cannot be used with `expiry_time_of_day` and `expiry_date`.
* `expiry_time_of_day` - **(Optional, String, HH:mm:ss)** The default time of day of expiration of entries within the cache.  Cannot be used with `expiry_timeout_in_sec` and `expiry_date`.
* `expiry_date` - **(Optional, String, MM-dd-yyyy)** The default date of expiration of entries within the cache.  Cannot be used with `expiry_timeout_in_sec` and `expiry_time_of_day`.
* `skip_cache_if_element_size_in_kb_exceeds` - **(Optional, Integer)** The maximum size of an entry in kilobytes that is allowed to be cached.
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Caches can be imported using a proper value of `id` as described above
