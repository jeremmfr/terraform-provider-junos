---
layout: "junos"
page_title: "Junos: junos_rib_group"
sidebar_current: "docs-junos-resource-rib-group"
description: |-
  Create a rib group
---

# junos_rib_group

Provides a rib group resource.

## Example Usage

```hcl
# Add a rib group
resource junos_rib_group "demo_rib" {
  name          = "prod"
  import_policy = ["policy-import-rib"]
  import_rib    = ["prod-vr.inet.0", "externe-vr.inet.0"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of rib group.
* `import_policy` - (Optional)(`ListOfString`) List of policy for import route.
* `import_rib` - (Optional)(`ListOfString`) List of import routing table
* `export_rib` - (Optional)(`String`) Export routing table

## Import

Junos rib group can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_rib_group.demo_rib prod
```
