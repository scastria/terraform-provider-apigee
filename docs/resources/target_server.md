---
subcategory: "Admin"
---
# Resource: apigee_target_server
Represents a named target server in an environment
## Example usage
```hcl
resource "apigee_target_server" "example" {
  environment_name = "dev"
  name = "Authentication"
  host = "auth.company.com"
  port = 80
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the target server
* `host` - **(Required, String)** The host name of the target server 
* `port` - **(Required, Integer)** The port of the target server
* `is_enabled` - **(Optional, Boolean)** Whether to enable this targer server for use
* `ssl_enabled` - **(Optional, Boolean)** Whether to communicate with this target server over TLS/SSL
* `ssl_keystore` - **(Optional, String)** Name of the keystore
* `ssl_keyalias` - **(Optional, String)** Name of the alias within the keystore
* `ssl_truststore` - **(Optional, String)** Name of the truststore that contains the certificate
* `ssl_client_auth_enabled` - **(Optional, Boolean)** Enable two-way TLS between Apigee and target
* `ssl_ignore_validation_errors` - **(Optional, Boolean)** Ignore TLS certificate errors
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
Target servers can be imported using a proper value of `id` as described above
