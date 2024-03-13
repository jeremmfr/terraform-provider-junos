---
page_title: "Junos: junos_lldp_interface"
---

# junos_lldp_interface

Provides a LLDP interface resource.

## Example Usage

```hcl
resource "junos_lldp_interface" "all" {
  name = "all"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Interface name or `all`.
- **disable** (Optional, Boolean)  
  Disable LLDP.
- **enable** (Optional, Boolean)  
  Enable LLDP.
- **power_negotiation** (Optional, Block)  
  LLDP power negotiation.
  - **disable** (Optional, Boolean)  
    Disable power negotiation.
  - **enable** (Optional, Boolean)  
    Enable power negotiation.
- **trap_notification_disable** (Optional, Boolean)  
  Disable lldp-trap notification.
- **trap_notification_enable** (Optional, Boolean)  
  Enable lldp-trap notification.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos lldp interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_lldp_interface.all all
```
