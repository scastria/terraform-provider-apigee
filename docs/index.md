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
//  access_token = "Use access token instead of username/password"
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
//  access_token = "Use access token instead of username/password"
  organization = "ZZZZ"
}

# Create a Role
resource "apigee_role" "example" {
  name = "WWWW"
}
```
## Argument Reference
* `username` - **(Optional, String)** The username that will invoke all Apigee API commands. Basic Authentication. Can be specified via env variable `APIGEE_USERNAME`.
* `password` - **(Optional, String)** The password of the username. Basic Authentication. Can be specified via env variable `APIGEE_PASSWORD`.
* `access_token` - **(Optional, String)** The access token from SAML or OAUTH authentication that can be used instead of `username` and `password`. Token Authentication. Can be specified via env variable `APIGEE_ACCESS_TOKEN`.
* `server` - **(Required, String)** The hostname of the Apigee Management API server. Default: `api.enterprise.apigee.com`. Can be specified via env variable `APIGEE_SERVER`.
* `private` - **(Required, Boolean)** Flag that determines whether the server is a private OnPrem installation of Apigee. Default: `false`. Can be specified via env variable `APIGEE_PRIVATE`.
* `port` - **(Required, Integer)** The port to use for the server. Default: `443`. Can be specified via env variable `APIGEE_PORT`.
* `organization` - **(Required, String)** The Apigee org that all Apigee API commands will work within. Can be specified via env variable `APIGEE_ORGANIZATION`.
