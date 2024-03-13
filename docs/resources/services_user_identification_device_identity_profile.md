---
page_title: "Junos: junos_services_user_identification_device_identity_profile"
---

# junos_services_user_identification_device_identity_profile

Provides a services user-identification device-information end-user-profile
(also named device identity profile) resource.

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

- **name** (Required, String, Forces new resource)  
  End-user-profile profile-name.
- **domain** (Required, String)  
  Domain name.
- **attribute** (Required, Block List)  
  For each name of attribute, configure list of values.
  - **name** (Required, String)  
    Attribute name.
  - **value** (Required, Set of String)  
    A list of values.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services user-identification device-information end-user-profile can be imported using an
id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_user_identification_device_identity_profile.demo demo
```
