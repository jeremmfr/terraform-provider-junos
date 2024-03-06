---
page_title: "Junos: junos_iccp"
---

# junos_iccp

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `protocols iccp` block.

Configure static configuration in `protocols iccp` block.

## Example Usage

```hcl
# Configure protocol ICCP
resource "junos_iccp" "iccp" {
  local_ip_addr                   = "192.0.2.1"
  authentication_key              = "a_key"
  session_establishment_hold_time = 300
}
```

## Argument Reference

The following arguments are supported:

- **local_ip_addr** (Required, String)  
  Local IP address to use by default for all peers.
- **authentication_key** (Optional, String, Sensitive)  
  MD5 authentication key for all peers.
- **session_establishment_hold_time** (Optional, Number)  
  Time within which connection must succeed with peers (45..600 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `iccp`.

## Import

Junos protocols ICCP can be imported using any id, e.g.

```shell
$ terraform import junos_iccp.iccp random
```
