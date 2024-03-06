---
page_title: "Junos: junos_snmp_community"
---

# junos_snmp_community

Provides a snmp community resource.

## Example Usage

```hcl
# Add a snmp community
resource "junos_snmp_community" "public" {
  name                    = "public"
  authorization_read_only = true
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of snmp community.
- **authorization_read_only** (Optional, Boolean)  
  Allow read-only access.  
  Conflict with `authorization_read_write`.
- **authorization_read_write** (Optional, Boolean)  
  Allow read and write access.  
  Conflict with `authorization_read_only`.
- **client_list_name** (Optional, String)  
  The name of client list or prefix list.  
  Conflict with `clients`.
- **clients** (Optional, Set of String)  
  List of source address prefix ranges to accept.  
  Conflict with `client_list_name`.
- **routing_instance** (Optional, Block List)  
  For each name of routing instance, accept clients.
  - **name** (Required, String)  
    Name of routing instance.
  - **client_list_name** (Optional, String)  
    The name of client list or prefix list.  
    Conflict with `clients` in block.
  - **clients** (Optional, Set of String)  
    List of source address prefix ranges to accept.  
    Conflict with `client_list_name` in block.
- **view** (Optional, String)  
  View name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos snmp community can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_snmp_community.public public
```
