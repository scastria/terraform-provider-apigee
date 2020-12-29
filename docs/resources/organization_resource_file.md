---
subcategory: "Develop"
---
# Resource: apigee_organization_resource_file
Represents a resource file in an organization
## Example usage
```hcl
resource "apigee_organization_resource_file" "example" {
  type = "js"
  name = "test.js"
  file = "resourceFiles/test.js"
  file_hash = filebase64sha256("resourceFiles/test.js")
}
```
## Argument Reference
* `type` - **(Required, ForceNew, String)** The type of the resource.
* `name` - **(Required, ForceNew, String)** The name of the resource.
* `file` - **(Required, String)** The filename of the resource.
* `file_hash` - **(Required, String)** The hash of the file used to detect changes of the contents of the file.
## Attribute Reference
* `id` - Same as `type`:`name`
## Import
Organization resource files can be imported using a proper value of `id` as described above
