---
subcategory: "Admin"
---
# Resource: apigee_virtual_host
Represents a virtual host in an environment
## Example usage
```hcl
resource "apigee_virtual_host" "example" {
  environment_name = "dev"
  name = "MainAPI"
  host_aliases = [
    "mainapi.company.com"
  ]
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the virtual host
* `host_aliases` - **(Required, List of String)** The aliases for the virtual host
* `port` - **(Optional, Integer)** The port of the virtual host
* `base_url` - **(Optional, String)** The base URL of the virtual host
* `ssl_enabled` - **(Optional, Boolean)** Whether to communicate with this virtual host over TLS/SSL
* `ssl_keystore` - **(Optional, String)** Name of the keystore
* `ssl_keyalias` - **(Optional, String)** Name of the alias within the keystore
* `ssl_truststore` - **(Optional, String)** Name of the truststore that contains the certificate
* `ssl_client_auth_enabled` - **(Optional, Boolean)** Enable two-way TLS between Apigee and target
* `ssl_ignore_validation_errors` - **(Optional, Boolean)** Ignore TLS certificate errors
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Virtual hosts can be imported using a proper value of `id` as described above
