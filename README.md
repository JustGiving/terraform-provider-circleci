# CircleCI Terraform provider

## JustGiving-specific changes

 - Added a `circleci_project` resource that follows a project on creation (enabling builds) and provides the project ID.
 - Added an example Terraform project using the provider in [`./example`](./example/).
 - Removed pipeline files

## Requirements

- [Terraform][terraform] 1.x
- [Go][go] 1.17+ (to build the provider plugin)

## Local testing

This will build the provider and output the binary in `~/.terraform.d/plugins/terraform.justgiving.com/justgiving/circleci/0.0.1/darwin_amd64/terraform-provider-circleci_v0.0.1`, where any Terraform stack on the local machine will be able to be loaded if required with:

```hcl
...
terraform {
  required_providers {
    circleci = {
      source = "terraform.justgiving.com/justgiving/circleci"
      version = "0.0.1"
    }
  }
}
```

## Publishing for service-generator

Build a new version for macOS and Linux (outputs the binaries in `./build`):

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags="-s -w" -a -o build/terraform-provider-circleci-darwin-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-s -w" -a -o build/terraform-provider-circleci-linux-amd64
```

Create a new release in the GitHub repository and upload the binaries:

```bash
version="x.y.z"
gh release create "v$version" -t '' -n '' build/*
```
