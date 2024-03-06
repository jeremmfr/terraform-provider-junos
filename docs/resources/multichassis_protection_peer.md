---
page_title: "Junos: junos_multichassis_protection_peer"
---

# junos_multichassis_protection_peer

Provides a multi-chassis inter-chassis protection peer resource.

## Example Usage

```hcl
# Add a multi-chassis inter-chassis protection peer
resource "junos_multichassis_protection_peer" "peer1" {
  ip_address = "192.0.2.1"
  interface  = "ge-0/0/3"
}
```

## Argument Reference

The following arguments are supported:

- **ip_address** (Required, String)  
  IP address for this peer.
- **interface** (Required, String)  
  Inter-Chassis protection link.
- **icl_down_delay** (Optional, Number)  
  Time in seconds between ICL down and MCAEs moving to standby (1..6000 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<ip_address>`.

## Import

Junos multi-chassis inter-chassis protection peer can be imported using an id made up of
`<ip_address>`, e.g.

```shell
$ terraform import junos_multichassis_protection_peer.peer1 192.0.2.1
```
