---
page_title: "Junos: junos_rib_group"
---

# junos_rib_group

Provides a rib group resource.

## Example Usage

```hcl
# Add a rib group
resource "junos_rib_group" "demo_rib" {
  name          = "prod"
  import_policy = ["policy-import-rib"]
  import_rib    = ["prod-vr.inet.0", "externe-vr.inet.0"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of rib group.
- **import_policy** (Optional, List of String)  
  List of policy for import route.
- **import_rib** (Optional, List of String)  
  List of import routing table
- **export_rib** (Optional, String)  
  Export routing table

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos rib group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_rib_group.demo_rib prod
```
