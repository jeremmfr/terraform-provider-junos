---
page_title: "Junos: junos_security_nat_source_pool"
---

# junos_security_nat_source_pool

Provides a security pool resource for source nat.

## Example Usage

```hcl
# Add a source nat pool
resource "junos_security_nat_source_pool" "demo_snat_pool" {
  name    = "ip_external"
  address = ["192.0.2.129/32"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Pool name.
- **address** (Required, List of String)  
  CIDR address to source nat pool.
- **address_pooling** (Optional, String)  
  Type of address pooling.  
  Need to be `paired` or `no-paired`.
- **description** (Optional, String)  
  Text description of pool.
- **pool_utilization_alarm_clear_threshold** (Optional, Number)  
  Lower threshold at which an SNMP trap is triggered.  
  Range 40 through 100.
- **pool_utilization_alarm_raise_threshold** (Optional, Number)  
  Upper threshold at which an SNMP trap is triggered.  
  Range 50 through 100.
- **port_no_translation** (Optional, Boolean)  
  Do not perform port translation.
- **port_overloading_factor** (Optional, Number)  
  Port overloading factor for each IP.
- **port_range** (Optional, String)  
  Range of port to source nat.  
  Format need to match `\d+(-\d+)?` with minimum low port 1024 and maximum high port 65535.
- **routing_instance** (Optional, String)  
  Name of routing instance to switch instance with nat.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security nat source pool can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_nat_source_pool.demo_snat_pool ip_external
```
