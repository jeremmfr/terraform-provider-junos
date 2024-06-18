---
page_title: "Junos: junos_switch_options"
---

# junos_switch_options

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `switch-options` block.  
By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `switch-options` block

## Example Usage

```hcl
# Configure switch-options
resource "junos_switch_options" "switch_options" {
  vtep_source_interface = "lo0.0"
}
```

## Argument Reference

The following arguments are supported:

- **clean_on_destroy** (Optional, Boolean)  
  Clean supported lines when destroy this resource.
- **remote_vtep_list** (Optional, Set of String)  
  Configure static remote VXLAN tunnel endpoints.
- **remote_vtep_v6_list** (Optional, Set of String)  
  Configure static ipv6 remote VXLAN tunnel endpoints.
- **service_id** (Optional, Number)  
  Service ID required if multi-chassis AE is part of a bridge-domain (1..65535).
- **vtep_source_interface** (Optional, String)  
  Source layer-3 IFL for VXLAN.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `switch_options`.

## Import

Junos switch_options can be imported using any id, e.g.

```shell
$ terraform import junos_switch_options.switch_options random
```
