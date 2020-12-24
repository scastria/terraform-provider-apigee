---
subcategory: "Develop"
---
# Resource: apigee_environment_resource_file
Represents a resource file in an environment
## Example usage
```hcl
resource "apigee_environment_resource_file" "example" {
  environment_name = "dev"
  type = "js"
  name = "test.js"
  file = "resourceFiles/test.js"
  file_hash = filebase64sha256("resourceFiles/test.js")
}
```
## Argument Reference
* `environment_name` - **(Required, ForceNew, String)** The name of an environment.
* `type` - **(Required, ForceNew, String)** The type of the resource.  Must be one of `java`, `js`, `jsc`, `hosted`, `node`, `py`, `wsdl`, `xsd`, or `xsl`.
* `name` - **(Required, ForceNew, String)** The name of the resource.
* `file` - **(Required, String)** The filename of the resource.
* `file_hash` - **(Required, String)** The hash of the file used to detect changes of the contents of the file.
## Attribute Reference
* `id` - Same as `environment_name`:`type`:`name`
## Import
Environment resource files can be imported using a proper value of `id` as described above
