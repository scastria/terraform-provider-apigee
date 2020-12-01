---
subcategory: "Develop"
---
# Resource: apigee_shared_flow
Represents a shared flow's latest revision
## Example usage
```hcl
resource "apigee_shared_flow" "example" {
  name = "MyFlow"
  bundle = "sharedflows/MyFlow/MyFlow.zip"
  bundle_hash = filebase64sha256("sharedflows/MyFlow/MyFlow.zip")
}
```
## Argument Reference
* `name` - **(Required, ForceNew, String)** The name of the shared flow.
* `bundle` - **(Required, String)** The filename of the bundle zip.
* `bundle_hash` - **(Required, String)** The hash of the bundle zip used to detect changes of the contents of the zip.
## Attribute Reference
* `id` - Same as `name`
* `revision` - The last revision imported
## Import
Shared flows can be imported using a proper value of `id` as described above
