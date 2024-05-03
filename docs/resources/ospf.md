---
page_title: "Junos: junos_ospf"
---

# junos_ospf

-> **Note:** This resource should only be created **once** for root level or each
routing-instance. It's used to configure static (not object) options in `protocols ospf` or
`protocols ospf3` block in root or routing-instance level.

Configure static configuration in `protocols ospf` or `protocols ospf3` block for root or
routing-instance level.

## Example Usage

```hcl
# Configure ospf
resource "junos_ospf" "ospf" {
  export = ["redistribe"]
}
```

## Argument Reference

The following arguments are supported:

- **version** (Optional, String)  
  Version of ospf.  
  Need to be `v2` or `v3`.  
  Defaults to `v2`.
- **routing_instance** (Optional, String, Forces new resource)  
  Routing instance.  
  Need to be `default` (for root level) or the name of routing instance.  
  Defaults to `default`.
- **database_protection** (Optional, Block)  
  Declare `database-protection` configuration.
  - **maximum_lsa** (Required, Number)  
    Maximum allowed non self-generated LSAs (1..1000000).
  - **ignore_count** (Optional, Number)  
    Maximum number of times to go into ignore state (1..32).
  - **ignore_time** (Optional, Number)  
    Time to stay in ignore state and ignore all neighbors (30..3600 seconds).
  - **reset_time** (Optional, Number)  
    Time after which the ignore count gets reset to zero (60..86400 seconds).
  - **warning_only** (Optional, Boolean)  
    Emit only a warning when LSA maximum limit is exceeded.
  - **warning_threshold** (Optional, Number)  
    Percentage of LSA maximum above which to trigger warning (30..100 percent).
- **disable** (Optional, Boolean)  
  Disable OSPF.
- **domain_id** (Optional, String)  
  Configure domain ID.  
  Only if `routing_instance` != `default`.
- **export** (Optional, List of String)  
  Export policy.
- **external_preference** (Optional, Number)  
  Preference of external routes.
- **forwarding_address_to_broadcast** (Optional, Boolean)  
  Set forwarding address in Type 5 LSA in broadcast network.
- **graceful_restart** (Optional, Block)  
  Declare `graceful-restart` configuration.
  - **disable** (Optional, Boolean)  
    Disable OSPF graceful restart capability.
  - **helper_disable** (Optional, Boolean)  
    Disable graceful restart helper capability.
  - **helper_disable_type** (Optional, String)  
    Disable graceful restart helper capability for specific type.  
    Need to be `both`, `restart-signaling` or `standard`.
  - **no_strict_lsa_checking** (Optional, Boolean)  
    Do not abort graceful helper mode upon LSA changes.
  - **notify_duration** (Optional, Number)  
    Time to send all max-aged grace LSAs (1..3600 seconds).
  - **restart_duration** (Optional, Number)  
    Time for all neighbors to become full (1..3600 seconds).
- **import** (Optional, List of String)  
  Import policy (for external routes or setting priority).
- **labeled_preference** (Optional, Number)  
  Preference of labeled routes.
- **lsa_refresh_interval** (Optional, Number)  
  LSA refresh interval (25..50 minutes).
- **no_nssa_abr** (Optional, Boolean)  
  Disable full NSSA functionality at ABR.
- **no_rfc1583** (Optional, Boolean)  
  Disable RFC1583 compatibility.
- **overload** (Optional, Block)  
  Set the overload mode (repel transit traffic).
  - **allow_route_leaking** (Optional, Boolean)  
    Allow routes to be leaked when overload is configured.
  - **as_external** (Optional, Boolean)  
    Advertise As External with maximum usable metric.
  - **stub_network** (Optional, Boolean)  
    Advertise Stub Network with maximum metric.
  - **timeout** (Optional, Number)  
    Time after which overload mode is reset (60..1800 seconds).
- **preference** (Optional, Number)  
  Preference of internal routes.
- **prefix_export_limit** (Optional, Number)  
  Maximum number of prefixes that can be exported (0..4294967295).
- **reference_bandwidth** (Optional, String)  
  Bandwidth for calculating metric defaults.
- **rib_group** (Optional, String)  
  Routing table group for importing OSPF routes.
- **sham_link** (Optional, Boolean)  
  Configure parameters for sham links.
- **sham_link_local** (Optional, String)  
  Local sham link endpoint address.
- **spf_options** (Optional, Block)  
  Declare `spf-options` configuration.
  - **delay** (Optional, Number)  
    Time to wait before running an SPF (50..8000 milliseconds).
  - **holddown** (Optional, Number)  
    Time to hold down before running an SPF (2000..20000 milliseconds).
  - **no_ignore_our_externals** (Optional, Boolean)  
    Do not ignore self-generated external and NSSA LSAs.
  - **rapid_runs** (Optional, Number)  
    Number of maximum rapid SPF runs before holddown (1..10).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<version>_-_<routing_instance>`.

## Import

Junos ospf can be imported using an id made up of `<version>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_ospf.ospf v2_-_default
```
