---
layout: "junos"
page_title: "Junos: junos_ospf"
sidebar_current: "docs-junos-resource-ospf"
description: |-
  Configure static configuration in protocols ospf or ospf3 block in root or routing-instance level.
---

# junos_ospf

-> **Note:** This resource should only be created **once** for root level or each routing-instance. It's used to configure static (not object) options in `protocols ospf` or `protocols ospf3` block in root or routing-instance level.

Configure static configuration in `protocols ospf` or `protocols ospf3` block for root ou routing-instance level.

## Example Usage

```hcl
# Configure ospf
resource junos_ospf "ospf" {
  export = ["redistribe"]
}
```

## Argument Reference

The following arguments are supported:

* `routing_instance` - (Optional)(`String`) Routing instance. Need to be 'default' (for root level) or the name of routing instance. Defaults to `default`.
* `version` - (Optional)(`String`) Version of ospf. Need to be 'v2' or 'v3'. Defaults to `v2`.
* `database_protection` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'database-protection' configuration.
  * `maximum_lsa` - (Required)(`Int`) Maximum allowed non self-generated LSAs (1..1000000).
  * `ignore_count` - (Optional)(`Int`) Maximum number of times to go into ignore state (1..32).
  * `ignore_time` - (Optional)(`Int`) Time to stay in ignore state and ignore all neighbors (30..3600 seconds).
  * `reset_time` - (Optional)(`Int`) Time after which the ignore count gets reset to zero (60..86400 seconds).
  * `warning_only` - (Optional)(`Bool`) Emit only a warning when LSA maximum limit is exceeded.
  * `warning_threshold` - (Optional)(`Int`) Percentage of LSA maximum above which to trigger warning (30..100 percent).
* `disable` - (Optional)(`Bool`) Disable OSPF.
* `domain_id` - (Optional)(`String`) Configure domain ID. Only if `routing_instance` != 'default'.
* `export` - (Optional)(`ListOfString`) Export policy.
* `external_preference` - (Optional)(`Int`) Preference of external routes.
* `forwarding_address_to_broadcast` - (Optional)(`Bool`) Set forwarding address in Type 5 LSA in broadcast network.
* `graceful_restart` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'graceful-restart' configuration.
  * `disable` - (Optional)(`Bool`) Disable OSPF graceful restart capability.
  * `helper_disable` - (Optional)(`Bool`) Disable graceful restart helper capability.
  * `helper_disable_type` - (Optional)(`String`) Disable graceful restart helper capability for specific type. Need to be 'both', 'restart-signaling' or 'standard'.
  * `no_strict_lsa_checking` - (Optional)(`Bool`) Do not abort graceful helper mode upon LSA changes.
  * `notify_duration` - (Optional)(`Int`) Time to send all max-aged grace LSAs (1..3600 seconds).
  * `restart_duration` - (Optional)(`Int`) Time for all neighbors to become full (1..3600 seconds).
* `import` - (Optional)(`ListOfString`) Import policy (for external routes or setting priority).
* `labeled_preference` - (Optional)(`Int`) Preference of labeled routes.
* `lsa_refresh_interval` - (Optional)(`Int`) LSA refresh interval (minutes) (25..50).
* `no_nssa_abr` - (Optional)(`Bool`) Disable full NSSA functionality at ABR.
* `no_rfc1583` - (Optional)(`Bool`) Disable RFC1583 compatibility.
* `overload` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to set the overload mode (repel transit traffic).
  * `allow_route_leaking` - (Optional)(`Bool`) Allow routes to be leaked when overload is configured.
  * `as_external` - (Optional)(`Bool`) Advertise As External with maximum usable metric.
  * `stub_network` - (Optional)(`Bool`) Advertise Stub Network with maximum metric.
  * `timeout` - (Optional)(`Int`) Time after which overload mode is reset (60..1800 seconds).
* `preference` - (Optional)(`Int`) Preference of internal routes.
* `prefix_export_limit` - (Optional)(`Int`) Maximum number of prefixes that can be exported (0..4294967295).
* `reference_bandwidth` - (Optional)(`String`) Bandwidth for calculating metric defaults.
* `rib_group` - (Optional)(`String`) Routing table group for importing OSPF routes.
* `sham_link` - (Optional)(`Bool`) Configure parameters for sham links.
* `sham_link_local` - (Optional)(`String`) Local sham link endpoint address.
* `spf_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'spf-options' configuration.
  * `delay` - (Optional)(`Int`) Time to wait before running an SPF (50..8000 milliseconds).
  * `holddown` - (Optional)(`Int`) Time to hold down before running an SPF (2000..20000 milliseconds).
  * `no_ignore_our_externals` - (Optional)(`Bool`) Do not ignore self-generated external and NSSA LSAs.
  * `rapid_runs` - (Optional)(`Int`) Number of maximum rapid SPF runs before holddown (1..10).

## Import

Junos ospf can be imported using an id made up of `<version>_-_<routing_instance>`, e.g.

```shell
$ terraform import junos_ospf.ospf v2_-_default
```
