---
page_title: "Junos: junos_oam_gretunnel_interface"
---

# junos_oam_gretunnel_interface

Provides a oam gre-tunnel interface resource.

## Example Usage

```hcl
# Add oam gre-tunnel interface
resource "junos_oam_gretunnel_interface" "gr1" {
  name           = "gr-1/1/10.1"
  hold_time      = 30
  keepalive_time = 10
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of interface.  
  Need to be a gr interface `gr-..`
- **hold_time** (Optional, Number)  
  Hold time (5..250 seconds).
- **keepalive_time** (Optional, Number)  
  Keepalive time (1..50 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos protocal oam gre-tunnel interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_oam_gretunnel_interface.gr1 gr-1/1/10.1
```
