# Apigee Provider
The Apigee provider is used to interact with the many resources supported by Apigee.  The provider needs to be configured with the proper credentials before it can be used.
## Example Usage
Terraform 0.13 and later:
```hcl
terraform {
  required_providers {
    apigee = {
      source  = "scastria/apigee"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Apigee Provider
provider "apigee" {
  username = "me@company.com"
  password = "XXXX"
  server = "api.enterprise.apigee.com"
  organization = "ZZZZ"
}

# Create a Role
resource "apigee_role" "example" {
  name = "WWWW"
}
```
Terraform 0.12 and earlier:
```hcl
# Configure the Apigee Provider
provider "apigee" {
  version = "~> 0.1.0"
  username = "me@company.com"
  password = "XXXX"
  server = "api.enterprise.apigee.com"
  organization = "ZZZZ"
}

# Create a Role
resource "apigee_role" "example" {
  name = "WWWW"
}
```
## Argument Reference
* `username` - **(Required, String)** The username that will invoke all Apigee API commands. Basic Authentication. Can be specified via env variable `APIGEE_USERNAME`.
* `password` - **(Required, String)** The password of the username. Basic Authentication. Can be specified via env variable `APIGEE_PASSWORD`.
* `server` - **(Required, String)** The hostname of the Apigee Management API server. Can be specified via env variable `APIGEE_SERVER`.
* `port` - **(Optional, String)** The port to use for the server. Default: 443. Can be specified via env variable `APIGEE_PORT`.
* `organization` - **(Required, String)** The Apigee org that all Apigee API commands will work within. Can be specified via env variable `APIGEE_ORGANIZATION`.
