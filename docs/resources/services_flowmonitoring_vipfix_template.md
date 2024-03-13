---
page_title: "Junos: junos_services_flowmonitoring_vipfix_template"
---

# junos_services_flowmonitoring_vipfix_template

Provides a services flow-monitoring version-ipfix template resource.

## Example Usage

```hcl
# Add a services flow-monitoring version-ipfix template
resource "junos_services_flowmonitoring_vipfix_template" "demo" {
  name = "demo"
  type = "ipv4-template"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of flow-monitoring version-ipfix template.
- **type** (Required, String)  
  Type of template.  
  Need to be `bridge-template`, `ipv4-template`, `ipv6-template` or `mpls-template`.
- **flow_active_timeout** (Optional, Number)  
  Interval after which active flow is exported (10..600).
- **flow_inactive_timeout** (Optional, Number)  
  Period of inactivity that marks a flow inactive (10..600).
- **flow_key_flow_direction** (Optional, Boolean)  
  Include flow direction.
- **flow_key_output_interface** (Optional, Boolean)  
  Include output interface.
- **flow_key_vlan_id** (Optional, Boolean)  
  Include vlan ID.
- **ip_template_export_extension** (Optional, Set of String)  
  Export-extension for `ipv4-template`, `ipv6-template` type.
- **mpls_template_label_position** (Optional, List of number)  
  One or more MPLS label positions (1..8).  
  `type` need to be `mpls-template`.
- **nexthop_learning_enable** (Optional, Boolean)  
  Enable nexthop learning.
- **nexthop_learning_disable** (Optional, Boolean)  
  Disable nexthop learning.
- **observation_domain_id** (Optional, Number)  
  Observation Domain Id (0..255).
- **option_refresh_rate** (Optional, Block)  
  Declare `option-refresh-rate` configuration.
  - **packets** (Optional, Number)  
    In number of packets (1..480000).
  - **seconds** (Optional, Number)  
    In number of seconds (10..600).
- **option_template_id** (Optional, Number)  
  Options template id (1024..65535).
- **template_id** (Optional, Number)  
  Template id (1024..65535).
- **template_refresh_rate** (Optional, Block)  
  Declare `template-refresh-rate` configuration.
  - **packets** (Optional, Number)  
    In number of packets (1..480000).
  - **seconds** (Optional, Number)  
    In number of seconds (10..600).
- **tunnel_observation_ipv4** (Optional, Boolean)  
  Tunnel observation IPv4.
- **tunnel_observation_ipv6** (Optional, Boolean)  
  Tunnel observation IPv6.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services flow-monitoring version-ipfix template can be imported using an id made up of
`<name>`, e.g.

```shell
$ terraform import junos_services_flowmonitoring_vipfix_template.demo demo
```
