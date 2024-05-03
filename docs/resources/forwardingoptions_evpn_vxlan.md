---
page_title: "Junos: junos_forwardingoptions_evpn_vxlan"
---

# junos_forwardingoptions_evpn_vxlan

-> **Note:** This resource should only be created **once** for root level or each routing-instance.
It's used to configure static (not object) options in `forwarding-options evpn-vxlan` block in root or
routing-instance level.

Configure static configuration in `forwarding-options evpn-vxlan` block for root or
routing-instance level.

## Example Usage

```hcl
# Configure forwarding-options evpn-vxlan
resource "junos_forwardingoptions_evpn_vxlan" "demo" {
  shared_tunnels = true
}
```

## Argument Reference

The following arguments are supported:

- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance if not root level.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **shared_tunnels** (Optional, Boolean)  
  Create VTEP tunnels to EVPN PE.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<routing_instance>`.

## Import

Junos forwarding-options evpn-vxlan can be imported using an id made up of
`<routing_instance>`, e.g.

```shell
$ terraform import junos_forwardingoptions_evpn_vxlan.demo default
```
