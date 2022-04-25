# Apigee Provider
The Apigee provider is used to interact with the many resources supported by Apigee.  The provider needs to be
configured with the proper credentials before it can be used.  The provider tries to work with all 3 known Apigee
hosting options:
1. [Apigee Edge for Private Cloud (OPDK)](https://apidocs.apigee.com/apis)
2. [Apigee Edge for Public Cloud (api.enterprise.apigee.com)](https://apidocs.apigee.com/apis)
3. [Apigee on Google Cloud (apigee.googleapis.com)](https://cloud.google.com/apigee/docs/reference)

However, each Apigee hosting option has different functionality so not all terraform resource types may be supported.
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
//  oauth_server = "Treat username/password as machine user and obtain access token automatically"
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
//  oauth_server = "Treat username/password as machine user and obtain access token automatically"
  organization = "ZZZZ"
}

# Create a Role
resource "apigee_role" "example" {
  name = "WWWW"
}
```
## Argument Reference
* `username` - **(Optional, String)** The username that will invoke all Apigee API commands. Basic Authentication. Or can be machine user for automatic machine user authentication using `oauth_server`. Can be specified via env variable `APIGEE_USERNAME`.
* `password` - **(Optional, String)** The password of the username. Basic Authentication. Or can be machine password for automatic machine user authentication using `oauth_server`. Can be specified via env variable `APIGEE_PASSWORD`.
* `access_token` - **(Optional, String)** The access token from SAML or OAUTH authentication that can be used instead of `username` and `password`. Token Authentication. Can be specified via env variable `APIGEE_ACCESS_TOKEN`.
* `use_ssl` - **(Optional, Boolean)** Whether to use https or http for Apigee server communication. Default: `true`.
* `server` - **(Optional, String)** The hostname of the Apigee Management API server. Default: `api.enterprise.apigee.com`. Can be specified via env variable `APIGEE_SERVER`.
* `server_path` - **(Optional, String)** The additional path of the Apigee Management API server. Default: `v1`. Can be specified via env variable `APIGEE_SERVER_PATH`.
* `port` - **(Optional, Integer)** The port to use for the server. Default: `443`. Can be specified via env variable `APIGEE_PORT`.
* `oauth_server` - **(Optional, String)** The hostname of the Apigee OAuth server that can generate access tokens for machine users. Can be specified via env variable `APIGEE_OAUTH_SERVER`.
* `oauth_server_path` - **(Optional, String)** The additional path of the Apigee OAuth server that can generate access tokens for machine users. Can be specified via env variable `APIGEE_OAUTH_SERVER_PATH`.
* `oauth_port` - **(Optional, Integer)** The port to use for the Apigee OAuth server. Default: `443`. Can be specified via env variable `APIGEE_OAUTH_PORT`.
* `organization` - **(Required, String)** The Apigee org that all Apigee API commands will work within. Can be specified via env variable `APIGEE_ORGANIZATION`.
