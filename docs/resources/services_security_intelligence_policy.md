---
page_title: "Junos: junos_services_security_intelligence_policy"
---

# junos_services_security_intelligence_policy

Provides a services security-intelligence policy resource.

## Example Usage

```hcl
# Add a services security-intelligence policy
resource "junos_services_security_intelligence_policy" "demo" {
  name = "demo"
  category {
    name         = "CC"
    profile_name = "profile_CC"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Security intelligence policy name.
- **category** (Required, Block List)  
  For each name of security intelligence category, configure a profile.
  - **name** (Required, String)  
    Name of security intelligence category.
  - **profile_name** (Required, String)  
    Name of profile.
- **description** (Optional, String)  
  Text description of policy.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services security-intelligence policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_security_intelligence_policy.demo demo
```
