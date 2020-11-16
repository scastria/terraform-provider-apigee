#Apigee Provider
The Apigee provider is used to interact with the many resources supported by Apigee.  The provider needs to be configured with the proper credentials before it can be used.
##Example Usage
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
  server = "YYYY"
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
  server = "YYYY"
  organization = "ZZZZ"
}

# Create a Role
resource "apigee_role" "example" {
  name = "WWWW"
}
```