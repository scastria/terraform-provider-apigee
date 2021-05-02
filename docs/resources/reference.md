---
subcategory: "Admin"
---
# Resource: apigee_reference
Represents a reference in an environment
## Example usage
```hcl
resource "apigee_keystore" "MyKeystore" {
  environment_name = "dev"
  name = "keystoreName"
}
resource "apigee_reference" "example" {
  environment_name = "dev"
  name = "refName"
  refers = apigee_keystore.MyKeystore.name
  resource_type = "KeyStore"
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment
* `name` - **(Required, ForceNew, String)** The name of the reference
* `refers` - **(Required, String)** Name of the keystore or truststore being referenced
* `resource_type` - **(Required, ForceNew, String)**  Set to KeyStore or TrustStore
## Attribute Reference
* `id` - Same as `environment_name`:`name`
## Import
References can be imported using a proper value of `id` as described above
