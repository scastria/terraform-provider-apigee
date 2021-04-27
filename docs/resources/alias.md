---
subcategory: "Admin"
---
# Resource: reference
Represents an alias in a keystore or truststore
## Example usage
```hcl
resource "testing_alias" "testingaliasthird" {
  environment_name = "dev"
  name = "alias-p12"
  keystore_name = "keystoreName" 
  format="pkcs12"
  file="certs/identity.p12"
  password="this is a test"
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** 
* `keystore_name` - **(Required, ForceNew, String)** The name of the keystore or truststore
* `format` - **(Required, ForceNew, String)** Type of alias creation. Valid values include: keycertjar, pkcs12, and  keycertfile.(selfsignedcert not supported yet)
* `file` - **(Optional, ForceNew, String)** The filename of a JAR or PKCS file
* `cert_file` - **(Optional, ForceNew, String)** The filename of a PEM file for the certFile 
* `key_file` - **(Optional, ForceNew, String)** The filename of a PEM file for the keyFile 
* `password` - **(Optional, ForceNew, String)** the password/passpharse in plain text 
* `password_env_var_name` - **(Optional, ForceNew, String)** name of the environment variable that holds the password/passphrase(this field overwrites password when present) 
* `ignore_expiry_validation` - **(Optional, String)** Flag that specifies whether to validate that the certificate hasn't expired. Set this value to true to skip validation. Defaults to false 
* `ignore_newline_validation` - **(Optional, String)** If false, do not throw an error when the file contains a chain with no newline between each cert. By default, Edge requires a newline between each cert in a chain. Defaults to true
## Attribute Reference
* `id` - Same as `keystore_name`:`name`
## Import
Alias can be imported using a proper value of `id` as described above
