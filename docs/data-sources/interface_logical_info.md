---
page_title: "Junos: junos_interface_logical_info"
---

# junos_interface_logical_info

Get summary information about a logical interface (like its admin/operational statuses and IP addresses).

## Example Usage

```hcl
# Read statuses and addresses of ge-0/0/3.0
data "junos_interface_logical_info" "ge003_0" {
  name = "ge-0/0/3.0"
}
output "ge003_0" {
  value = data.junos_interface_logical_info.ge003_0
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name of unit interface (with dot).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  The name of interface read.
- **admin_status** (String)  
  Admin status.
- **oper_status** (String)  
  Operational status.
- **family_inet** (Block)  
  Family inet enabled.
  - **address_cidr** (List of String)  
    List of addresses in CIDR format.
- **family_inet6** (Block)  
  Family inet6 enabled.
  - **address_cidr** (List of String)  
    List of addresses in CIDR format.
