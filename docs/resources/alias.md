---
subcategory: "Admin"
---
# Resource: apigee_alias
Represents a keystore alias in an environment
## Example usage
```hcl
resource "apigee_keystore" "MyKeystore" {
  environment_name = "dev"
  name = "MyKeystore"
}
resource "apigee_alias" "example" {
  environment_name = "dev"
  keystore_name = apigee_keystore.MyKeystore.name
  name = "MyAlias"
  format = "pkcs12"
  file = "cert.p12"
  file_hash = filebase64sha256("cert.p12")
  password = "certpassword"
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `keystore_name` - **(Required, ForceNew, String)** The name of a keystore
* `name` - **(Required, ForceNew, String)** The name of the alias
* `format` - **(Required, String)** The format of input files used to upload alias cert/key. Allowed values: `keycertfile`, `keycertjar`, and `pkcs12`.
* `file` - **(Optional, String)** The filename used for formats: `keycertjar` and `pkcs12`.
* `file_hash` - **(Optional, String)** The hash of the `file` used to detect changes of the contents of the `file`.
* `key_file` - **(Optional, String)** The key filename used for format: `keycertfile`.
* `key_file_hash` - **(Optional, String)** The hash of the `key_file` used to detect changes of the contents of the `key_file`.
* `cert_file` - **(Optional, String)** The cert filename used for format: `keycertfile`.
* `cert_file_hash` - **(Optional, String)** The hash of the `cert_file` used to detect changes of the contents of the `cert_file`.
* `password` - **(Optional, String)** The password of any file containing a key.
* `ignore_expiry_validation` - **(Optional, Boolean)** Flag that specifies whether to validate that the certificate hasn't expired. Set this value to `true` to skip validation.
* `ignore_newline_validation` - **(Optional, Boolean)** If `false`, do not throw an error when the file contains a chain with no newline between each cert.
## Attribute Reference
* `id` - Same as `environment_name`:`keystore_name`:`name`
## Import
Aliases can be imported using a proper value of `id` as described above.  Apigee does not allow determining the original format used to initially create an alias.  Therefore, importing will not result in a complete state.
## Updating
Apigee only allows updating a certificate of an existing alias via the `file` property.
