---
page_title: "Junos: junos_forwardingoptions_dhcprelay_servergroup"
---

# junos_forwardingoptions_dhcprelay_servergroup

Provides a DHCP relay server group.

## Example Usage

```hcl
# Add a DHCP relay server group
resource "junos_forwardingoptions_dhcprelay_servergroup" "demo" {
  name = "demo"
  ip_address = [
    "192.0.2.8",
  ]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Server group name.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance for server group.  
  Need to be `default` or name of routing instance.  
  Defaults to `default`
- **version** (Optional, String, Forces new resource)  
  Version for DHCP or DHCPv6.  
  Need to be `v4` or `v6`.
- **ip_address** (Optional, List of String)  
  IP Addresses of DHCP servers.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>_-_<routing_instance>_-_<version>`.

## Import

Junos forwarding-options dhcp-relay server-group can be imported using an id made up of
`<name>_-_<routing_instance>_-_<version>`, e.g.

```shell
$ terraform import junos_forwardingoptions_dhcprelay_servergroup.demo demo_-_default_-_v4
```
