---
layout: "junos"
page_title: "Junos: junos_services_user_identification_device_identity_profile"
sidebar_current: "docs-junos-resource-services-user-identification-device-identity-profile"
description: |-
  Create a services user-identification device-information end-user-profile (also named device identity profile)
---

# junos_services_user_identification_device_identity_profile

Provides a services user-identification device-information end-user-profile (also named device identity profile) resource.

## Example Usage

```hcl
# Add a services user-identification device-information end-user-profile
resource "junos_services_user_identification_device_identity_profile" "demo" {
  name   = "demo"
  domain = "domain"
  attribute {
    name  = "device-identity"
    value = ["device1", "barcode scan"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) End-user-profile profile-name.
* `domain` - (Required)(`String`) Domain name.
* `attribute` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure a attribute. Can be specified multiple times for each attribute.
  * `name` - (Required)(`String`) Attribute name.
  * `value` - (Required)(`ListOfString`) A list of values.

## Import

Junos services user-identification device-information end-user-profile can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_services_user_identification_device_identity_profile.demo demo
```
