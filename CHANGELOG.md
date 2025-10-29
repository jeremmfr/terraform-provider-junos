<!-- markdownlint-disable-file MD013 MD041 -->
# changelog

## v2.15.0 (2025-10-29)

ENHANCEMENTS:

* release now with Go 1.25
* **resource/junos_security_global_policy**: add `description` argument in `policy` block
* **resource/junos_security_policy**: add `description` argument in `policy` block (Fix [#834](https://github.com/jeremmfr/terraform-provider-junos/issues/834))

## v2.14.0 (2025-07-30)

FEATURES:

* add `junos_security_utm_custom_message` resource

ENHANCEMENTS:

* **resource/junos_chassis_cluster**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_chassis_redundancy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_null_commit_file**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  `triggers` argument now accept any attribute type (and so still a Map)  
  the permissions of file targeted by `filename` argument are now preserved when using `clear_file_after_commit` argument
* **resource/junos_rib_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value
* **resource/junos_security_utm_custom_url_category**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_custom_url_pattern**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `content_filtering_rule_set` block argument
* **resource/junos_security_utm_profile_web_filtering_juniper_enhanced**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `custom_message` argument in root level and `category` block
* **resource/junos_security_utm_profile_web_filtering_juniper_local**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `no_safe_search` argument
  * add `custom_message` argument
  * add `category` block argument
* **resource/junos_security_utm_profile_web_filtering_websense_redirect**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `no_safe_search` argument
  * add `routing_instance` and `source_address` arguments in `server` block
  * add `custom_message` argument
  * add `category` block argument
* **resource/junos_services_advanced_anti_malware_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_services_proxy_profile**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* the provider don't use anymore the legacy SDKv2 plugin and the mux plugin used during migration to plugin framework

## v2.13.0 (2025-06-05)

ENHANCEMENTS:

* release now with golang 1.24
* **resource/junos_layer2_control**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_lldp_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_lldpmed_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_services_security_intelligence_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value
* **resource/junos_services_security_intelligence_profile**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `action` argument in `default_rule_then` block and `then_action` argument in `rule` block accept now `sinkhole` for `DNS` category
* **resource/junos_services**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list  
  computed attributes, in `advanced_anti_malware.connection` and `security_intelligence` block, are now unknown if block is specified without these attributes in Plan to update resources (except if attributes are null in Terraform state), instead of using values in Terraform state
* **resource/junos_services_ssl_initiation_profile**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_services_user_identification_ad_access_domain**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_services_user_identification_device_identity_profile**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
* **resource/junos_vlan**: add `interface` argument (Fix [#794](https://github.com/jeremmfr/terraform-provider-junos/issues/794))

## v2.12.0 (2025-02-24)

ENHANCEMENTS:

* **resource/junos_forwardingoptions_dhcprelay_servergroup**: `ip_address` argument is no longer ordered
* **resource/junos_igmp_snooping_vlan**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_interface_logical**:
  * add `virtual_gateway_accept_data`, `virtual_gateway_v4_mac` and `virtual_gateway_v6_mac` arguments
  * add `virtual_gateway_address` argument in `address` block in `family_inet` and `family_inet6` blocks
  * move config errors from defined `vrrp_group` block on `st0.` interface during Plan instead of during Apply
* **data-source/junos_interface_logical**: add `virtual_gateway_accept_data`, `virtual_gateway_v4_mac`, `virtual_gateway_v6_mac` and `virtual_gateway_address` attributes like resource
* **resource/junos_routing_options**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `ipv6_router_id` argument
* **resource/junos_security**: add `routing_instance` argument in `idp_security_package` block (Fix [#754](https://github.com/jeremmfr/terraform-provider-junos/issues/754))
* **resource/junos_system_ntp_server**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attribute doesn't accept value *false*  
  optional string attribute doesn't accept *empty* value  
  * add `nts` block argument
* **resource/junos_security_dynamic_address_feed_server**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_security_dynamic_address_name**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `session_scan` argument
* **resource/junos_system_services_dhcp_localserver_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **data-source/junos_system_information** data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)

BUG FIXES:

* **resource/junos_forwardingoptions_dhcprelay**:
  * fix missing detection of conflict between `dynamic_profile_aggregate_clients` and `dynamic_profile_use_primary` arguments in config validation
  * fix missing detection of empty `overrides_v4` and `overrides_v6` block arguments in config validation
* **resource/junos_forwardingoptions_dhcprelay_group**:
  * fix missing detection of conflict between `dynamic_profile_aggregate_clients` and `dynamic_profile_use_primary` arguments in config validation
  * fix missing detection of empty `overrides_v4` and `overrides_v6` block arguments in config validation

## v2.11.0 (2025-01-08)

FEATURES:

* add **junos_mstp** resource (Partial fix [#732](https://github.com/jeremmfr/terraform-provider-junos/issues/732))
* add **junos_mstp_interface** resource (Partial fix [#732](https://github.com/jeremmfr/terraform-provider-junos/issues/732))
* add **junos_mstp_msti** resource (Partial fix [#732](https://github.com/jeremmfr/terraform-provider-junos/issues/732))

ENHANCEMENTS:

* **resource/***: set and list attributes, in resources using the new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework), now test that they contain only non-null values
* **resource/junos_access_address_assignment_pool**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `user_name` argument inside `host` block inside `family` block (`hardware_address` is now optional)
  * `host` block, `xauth_attributes_primary_dns` and `xauth_attributes_secondary_dns` arguments inside `family` block can now be configured when `type` = `inet6`
* **resource/junos_group_dual_system**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false* (but still accepted for `apply_groups` argument)  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_idp_custom_attack**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `recommended_action` argument is now optional
  * values of `ip_flags` argument in `protocol_ipv4` block has now a plan validator to one of `df`, `mf`, `rb`, `no-df`, `no-mf` or `no-rb`
  * values of `tcp_flags` argument in `protocol_tcp` block has now a plan validator to one of `ack`, `fin`, `psh`, `r1`, `r2`, `rst`, `syn`, `urg`, `no-acl`, `no-fin`, `no-psh`, `no-r1`, `no-r2`, `no-rst`, `no-syn` or `no-urg`
* **resource/junos_security_idp_custom_attack_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_security_idp_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_screen**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `aggregation` block argument
* **resource/junos_security_screen_whitelist**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_services_rpm_probe**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_system_login_class**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  * values of `allowed_days` argument now have a plan validator to be one of the valid days (`sunday`, `monday`, ...)  
  * values of `permissions` argument now have a plan validator to be one of the valid permission flags (`access`, `access-control`, `admin`, ...)  
* **resource/junos_system_login_user**, **resource/junos_system_root_authentication**: mark `encrypted_password` attribute as sensitive (to obscure the value in Terraform output) even if the data is encrypted

BUG FIXES:

* **resource/junos_forwardingoptions_dhcprelay**: fix validation of `attempts` must be in range (1..10) when `version` = v6 in `bulk_leasequery` block

## v2.10.0 (2024-10-25)

FEATURES:

* add **junos_applications_ordered** resource, copy of **junos_applications** resource but with Block List instead of Block Set to have a workaround for the performance issue on Block Sets (workaround for [#709](https://github.com/jeremmfr/terraform-provider-junos/issues/709))
* add **junos_security_address_book_ordered** resource, copy of **junos_security_address_book** resource but with Block List instead of Block Set to have a workaround for the performance issue on Block Sets (workaround for [#498](https://github.com/jeremmfr/terraform-provider-junos/issues/498))
* add **junos_security_global_policy_unordered** resource, copy of **junos_security_global_policy** resource but with Block Set instead of Block List to have a workaround for too complex plan output when the number of blocks on the resource changes
* add **junos_security_policy_unordered** resource, copy of **junos_security_policy** resource but with Block Set instead of Block List to have a workaround for too complex plan output when the number of blocks on the resource changes
* add **junos_security_zone_ordered** resource, copy of **junos_security_zone** resource but with Block List instead of Block Set to have a workaround for the performance issue on Block Sets

ENHANCEMENTS:

* **resource/junos_forwardingoptions_dhcprelay**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_forwardingoptions_dhcprelay_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_forwardingoptions_dhcprelay_servergroup**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)

BUG FIXES:

* **resource/***: remove the action of delete candidate configuration before unlocking it, as this is unnecessary and may produce inconsistent warnings (Fix [#710](https://github.com/jeremmfr/terraform-provider-junos/issues/710))
* **resource/junos_bridge_domain**: fix missing validate of not empty resource in create/update functions
* **resource/junos_forwardingoptions_evpn_vxlan**: fix missing validate of not empty resource in create/update functions
* **resource/junos_ospf_area**: fix missing part of validate config when `version` is null in config (default to `v2`)
* **resource/junos_vlan**: fix missing validate of not empty resource in create/update functions

## v2.9.0 (2024-09-10)

FEATURES:

* add **junos_applications** resource (Fix [#694](https://github.com/jeremmfr/terraform-provider-junos/issues/694))
* add **junos_security_authentication_key_chain** resource
* **provider**: add `no_decode_secrets` attribute to disable decoding secret `$9$` hashes by Junos device when reading resource data (Fix [#688](https://github.com/jeremmfr/terraform-provider-junos/issues/688))

ENHANCEMENTS:

* **resource/junos_application**: add `do_not_translate_a_query_to_aaaa_query`, `do_not_translate_aaaa_query_to_a_query`, `icmp_code`, `icmp_type`, `icmp6_code` and `icmp6_type` arguments
* **data-source/junos_applications**:
  * add `do_not_translate_a_query_to_aaaa_query`, `do_not_translate_aaaa_query_to_a_query`, `icmp_code`, `icmp_type`, `icmp6_code` and `icmp6_type` attributes in `applications` attribute
  * add `do_not_translate_a_query_to_aaaa_query` and `do_not_translate_aaaa_query_to_a_query` arguments inside `match_options` block argument
* **resource/junos_rip_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_rip_neighbor**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_snmp_v3_usm_user**: now provider store the corresponding `authentication-key` and `privacy-key` of `authentication-password` and `privacy-password` in private state of Terraform after create/update resource to be able to detect a change of the password outside of Terraform.
* **resource/junos_system**: add `keys` argument inside `license` block argument (Fix [#689](https://github.com/jeremmfr/terraform-provider-junos/issues/689))
* **resource/junos_system_login_user**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * now provider store the corresponding `authentication encrypted-password` of `authentication plain-text-password` in private state of Terraform after create/update resource to be able to detect a change of the password outside of Terraform.
* **resource/junos_system_root_authentication**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  * now provider store the corresponding `encrypted-password` of `plain-text-password` in private state of Terraform after create/update resource to be able to detect a change of the password outside of Terraform.
* release now with golang 1.23

BUG FIXES:

* **resource/junos_security_ike_policy**, **resource/junos_security_ike_proposal**, **resource/junos_security_ipsec_policy**, **resource/junos_security_ipsec_proposal**, **resource/junos_security_ipsec_vpn**: don't check device compatibility with security model (could be used on non-security devices)

## v2.8.0 (2024-06-27)

ENHANCEMENTS:

* **resource/junos_bridge_domain**: add `static_remote_vtep_list` argument inside `vxlan` block argument (Fix [#672](https://github.com/jeremmfr/terraform-provider-junos/issues/672))
* **resource/junos_interface_logical**: add `encapsulation` argument (Fix [#674](https://github.com/jeremmfr/terraform-provider-junos/issues/674))
* **data-source/junos_interface_logical**: add `encapsulation` attribute like resource
* **resource/junos_routing_instance**: add `remote_vtep_list` and `remote_vtep_v6_list` arguments (Fix [#673](https://github.com/jeremmfr/terraform-provider-junos/issues/673))
* **data-source/junos_routing_instance**: add `remote_vtep_list` and `remote_vtep_v6_list` attributes like resource
* **resource/junos_rstp**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_rstp_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_security_log_stream**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `transport` block argument (Fix [#675](https://github.com/jeremmfr/terraform-provider-junos/issues/675))
* **resource/junos_switch_options**: add `remote_vtep_list` and `remote_vtep_v6_list` arguments
* **resource/junos_vlan**: add `static_remote_vtep_list` argument inside `vxlan` block argument
* **resource/junos_vstp**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_vstp_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_vstp_vlan**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
* **resource/junos_vstp_vlan_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*

## v2.7.0 (2024-05-03)

FEATURES:

* add **junos_forwardingoptions_evpn_vxlan** resource (Partial Fix [#645](https://github.com/jeremmfr/terraform-provider-junos/issues/645))

ENHANCEMENTS:

* **data-source/junos_interfaces_physical_present**:
  * add `interfaces` block map attribute with same attributes as `interface_statuses` and additional `logical_interface_names` attribute (Fix [#641](https://github.com/jeremmfr/terraform-provider-junos/issues/641))
  * deprecate `interface_statuses` attribute (read the `interfaces` attribute instead)
* **resource/junos_evpn**: add `no_core_isolation` argument (Fix [#644](https://github.com/jeremmfr/terraform-provider-junos/issues/644))
* **resource/junos_ospf_area**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_ospf**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_vlan**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `community_vlans` argument is now a Set of String (instead of Set of Number) to accept VLAN name in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform
  * `isolated_vlan` argument is now a String (instead of Number) to accept VLAN name in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform
  * `vlan_id` argument is now a String (instead of Number) to accept `all` or `none` in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform
  * add `routing_instance` argument (Partial fix [#646](https://github.com/jeremmfr/terraform-provider-junos/issues/646))  
    and therefore `id` format has been changed to `<name>_-_<routing_instance>` (instead of `<name>`)
  * add `no_arp_suppression` argument (Partial fix [#646](https://github.com/jeremmfr/terraform-provider-junos/issues/646))
  * add `translation_vni` argument inside `vxlan` block argument

## v2.6.0 (2024-03-13)

FEATURES:

* add **junos_system_tacplus_server** resource (Fix [#629](https://github.com/jeremmfr/terraform-provider-junos/issues/629))
* add **junos_virtual_chassis** resource (Fix [#623](https://github.com/jeremmfr/terraform-provider-junos/issues/623))

ENHANCEMENTS:

* **resource/junos_eventoptions_destination**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_eventoptions_generate_event**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    optional boolean attributes doesn't accept value *false*  
  * add `start_time` argument
* **resource/junos_eventoptions_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_snmp**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_snmp_clientlist**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_snmp_community**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_snmp_v3_community**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_snmp_v3_usm_user**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  * `authentication_type` argument accept new value: `authentication-sha224`, `authentication-sha256`, `authentication-sha384` and `authentication-sha512`
* **resource/junos_snmp_v3_vacm_accessgroup**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_snmp_v3_vacm_securitytogroup**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
* **resource/junos_snmp_view**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_system**:
  * add `accounting` block argument (Fix [#630](https://github.com/jeremmfr/terraform-provider-junos/issues/630))
  * add `radius_options_attributes_nas_id` argument
  * add `tacplus_options_authorization_time_interval`, `tacplus_options_enhanced_accounting`, `tacplus_options_exclude_cmd_attribute`, `tacplus_options_no_cmd_attribute_value`, `tacplus_options_service_name`, `tacplus_options_strict_authorization`, `tacplus_options_no_strict_authorization`, `tacplus_options_timestamp_and_timezone` arguments
* **resource/junos_system_radius_server**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
* release now with golang 1.22

## v2.5.0 (2024-01-25)

FEATURES:

* add **junos_chassis_inventory** data-source (Fix [#587](https://github.com/jeremmfr/terraform-provider-junos/issues/587))

BUG FIXES:

* **resource/junos_aggregate_route, junos_application, junos_bgp_group, junos_bgp_neighbor, junos_bridge_domain, junos_evpn, junos_firewall_filter, junos_firewall_policer, junos_forwardingoptions_sampling_instance, junos_forwardingoptions_sampling, junos_forwardingoptions_storm_control_profile, junos_generate_route, junos_interface_logical, junos_interface_physical, junos_policyoptions_community, junos_routing_instance, junos_security, junos_security_address_book, junos_security_global_policy, junos_security_ike_gateway, junos_security_ike_policy, junos_security_ipsec_policy, junos_security_ipsec_vpn, junos_security_nat_destination_pool, junos_security_nat_destination, junos_security_nat_source_pool, junos_security_nat_static_rule, junos_security_nat_static, junos_security_policy, junos_security_zone_book_address, junos_security_zone, junos_services_flowmonitoring_v9_template, junos_services_flowmonitoring_vipfix_template, junos_static_route, junos_system_syslog_file, junos_system**:  
avoid triggering the conflict errors when Terraform calls the resource config validate function and the value for potential conflict is unknown (can be null afterward) (Fix [#611](https://github.com/jeremmfr/terraform-provider-junos/issues/611))

## v2.4.0 (2023-12-21)

FEATURES:

* add **junos_forwardingoptions_storm_control_profile** resource (Partial fix [#574](https://github.com/jeremmfr/terraform-provider-junos/issues/574))
* add **junos_iccp** resource (Partial fix [#573](https://github.com/jeremmfr/terraform-provider-junos/issues/573))
* add **junos_iccp_peer** resource (Partial fix [#573](https://github.com/jeremmfr/terraform-provider-junos/issues/573))
* add **junos_multichassis** resource (Partial fix [#576](https://github.com/jeremmfr/terraform-provider-junos/issues/576))
* add **junos_multichassis_protection_peer** resource (Partial fix [#576](https://github.com/jeremmfr/terraform-provider-junos/issues/576))
* add **junos_system_syslog_user** resource (Fix [#593](https://github.com/jeremmfr/terraform-provider-junos/issues/593))
* **provider**: add `commit_confirmed` and `commit_confirmed_wait_percent` argument to be able use `commit confirmed` feature to commit the resource actions (Fix [#585](https://github.com/jeremmfr/terraform-provider-junos/issues/585))

ENHANCEMENTS:

* **resource/junos_interface_physical**:
  * add `mc_ae` block argument in `parent_ether_opts` block (Fix [#572](https://github.com/jeremmfr/terraform-provider-junos/issues/572))
  * add `storm_control` argument (Partial fix [#574](https://github.com/jeremmfr/terraform-provider-junos/issues/574))
* **data-source/junos_interface_physical**:
  * add `mc_ae` block attribute in `parent_ether_opts` block like resource
  * add `storm_control` attribute like resource
* **resource/junos_switch_options**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
  * add `service_id` argument (Fix [#575](https://github.com/jeremmfr/terraform-provider-junos/issues/575))
* **resource/junos_system**: add `web_management_session_idle_timeout` and `web_management_session_limit` arguments in `services` block (Fix [#594](https://github.com/jeremmfr/terraform-provider-junos/issues/594))
* **resource/junos_system_syslog_file**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_system_syslog_host**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **provider**: display all errors when configuration commit generate multiple errors

BUG FIXES:

* **data-source/junos_interface_physical**: fix reading `link_speed` and `minimum_bandwidth` attributes in `parent_ether_opts` block
* **resource/junos_system_syslog_file**: fix reading `archive size` when value is a multiple of 1024 (k,m,g)

## v2.3.3 (2023-12-07)

BUG FIXES:

* **resource/junos_system**: fix crash when `web_management_https` is defined and `web_management_http` is not (Fix [#588](https://github.com/jeremmfr/terraform-provider-junos/issues/588))

## v2.3.2 (2023-11-16)

BUG FIXES:

* **resource/junos_firewall_filter**: allow use `protocol` and `protocol_except` arguments in `from` block in `term` block when `family` is set to `ethernet-switching` (Fix [#577](https://github.com/jeremmfr/terraform-provider-junos/issues/577))

## v2.3.1 (2023-11-10)

BUG FIXES:

* **resource/junos_system**: fix value validator (also accept `@`, `.`) on `ciphers`, `hostkey_algorithm`, `key_exchange` and `macs` attributes in `ssh` block in `services` block (Fix [#570](https://github.com/jeremmfr/terraform-provider-junos/issues/570))

## v2.3.0 (2023-11-08)

ENHANCEMENTS:

* **resource/junos_bridge_domain**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `interface` argument (Fix [#548](https://github.com/jeremmfr/terraform-provider-junos/issues/548))
* **resource/junos_evpn**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `duplicate_mac_detection` block argument (Fix [#535](https://github.com/jeremmfr/terraform-provider-junos/issues/535))
* **resource/junos_system**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `authentication_order`, `auxiliary_authentication_order`, `console_authentication_order` arguments have now a value validator: need to be `password`, `radius` or `tacplus`
  * add `name_server_opts` argument (in conflict with `name_server` argument) to also configure DNS name server but with optional options (`routing_instance`) (Fix [#561](https://github.com/jeremmfr/terraform-provider-junos/issues/561))

BUG FIXES:

* **resource/junos_aggregate_route**, **resource/junos_generate_route**, **resource/junos_static_route**: fix missing no-empty value validator on `as_path_path` and `next_table` arguments

## v2.2.0 (2023-09-13)

ENHANCEMENTS:

* **resource/junos_interface_physical**: add `trunk_non_els` and `vlan_native_non_els` arguments (Fix [#521](https://github.com/jeremmfr/terraform-provider-junos/issues/521))
* **data-source/junos_interface_physical**: add `trunk_non_els` and `vlan_native_non_els` attributes
* **resource/junos_aggregate_route**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value
* **resource/junos_generate_route**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value
* **resource/junos_static_route**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value
* **data-source/junos_routes**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_nat_destination**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_nat_destination_pool**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
* **resource/junos_security_nat_source**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_nat_source_pool**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
* **resource/junos_security_nat_static**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_nat_static_rule**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* release now with golang 1.21

## v2.1.3 (2023-08-30)

BUG FIXES:

* **resource/junos_forwardingoptions_sampling_instance**: avoid resources replacement when upgrading the provider before `v2.0.0` and without refreshing resource states (`-refresh=false`) (Fix [#536](https://github.com/jeremmfr/terraform-provider-junos/issues/536))

## v2.1.2 (2023-08-28)

BUG FIXES:

* **resource/junos_security_ike_gateway** : fix `Value Conversion Error` when upgrading the provider before `v2.0.0` to `v2.0.0...v2.1.1` and there are this type of resource with `remote_identity` block set in state (Fix [#533](https://github.com/jeremmfr/terraform-provider-junos/issues/533))

## v2.1.1 (2023-08-21)

BUG FIXES:

* **resource/junos_policyoptions_policy_statement**: fix potential crash when applying (create/update) resource (Fix [#528](https://github.com/jeremmfr/terraform-provider-junos/issues/528))

## v2.1.0 (2023-07-25)

ENHANCEMENTS:

* **resource/junos_application**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value
* **data-source/junos_applications**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_application_set**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  * add `application_set`, `description` arguments
* **data-source/junos_application_sets**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
  * add `match_application_sets` argument
  * add `application_set` and `description` attribute in `application_sets` block attribute
* **resource/junos_bgp_group**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `advertise_external` is now computed to `true` when `advertise_external_conditional` is `true` (instead of the 2 attributes conflicting)
  * `bfd_liveness_detection.version` now generate an error if the value is not in one of strings `0`, `1` or `automatic`
  * add `bgp_error_tolerance`, `description`, `no_client_reflect`, `tcp_aggressive_transmission` arguments
* **resource/junos_bgp_neighbor**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `advertise_external` is now computed to `true` when `advertise_external_conditional` is `true` (instead of the 2 attributes conflicting)
  * `bfd_liveness_detection.version` now generate an error if the value is not in one of strings `0`, `1` or `automatic`
  * add `bgp_error_tolerance`, `description`, `no_client_reflect`, `tcp_aggressive_transmission` arguments
* **resource/junos_firewall_filter**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `destination_mac_address`, `destination_mac_address_except`,
  `forwarding_class`, `forwarding_class_except`,
  `interface`,
  `loss_priority`, `loss_priority_except`,
  `packet_length`, `packet_length_except`,
  `policy_map`, `policy_map_except`,
  `source_mac_address` and `source_mac_address_except` arguments in from block in term block
  * add `forwarding_class` and `loss_priority` arguments in then block in term block
* **resource/junos_firewall_policer**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `logical_bandwidth_policer`, `logical_interface_policer`, `physical_interface_policer`, `shared_bandwidth_policer` and `if_exceeding_pps` arguments
* **resource/junos_policyoptions_as_path**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
* **resource/junos_policyoptions_as_path_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
* **resource/junos_policyoptions_community**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    optional boolean attributes doesn't accept value *false*  
  * add `dynamic_db` argument (`members` is now optional but one of `dynamic_db` or `members` must be specified)
* **resource/junos_policyoptions_policy_statement**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `dynamic_db` argument
* **resource/junos_policyoptions_prefix_list**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  

BUG FIXES:

* reduce plan time for resources that have migrated to the new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and have block set attributes (multiple unordered blocks) (Partial fix [#498](https://github.com/jeremmfr/terraform-provider-junos/issues/498))
* **resource/junos_security_ipsec_vpn**: fix length validator (max 31 instead of 32) and remove space exclusion validator of `name` for `traffic_selector` block

## v2.0.0 (2023-05-10)

BREAKING CHANGES with new `v2`:

* **provider**: remove `aes128-cbc` cipher from default ciphers when `ssh_ciphers` is not specified
* remove deprecated **junos_interface** resource and data source
* **resource/junos_bgp_group**, **resource/junos_bgp_neighbor**: remove deprecated `multipath` argument
* **resource/junos_interface_physical**: remove deprecated `ae_lacp`, `ae_link_speed`, `ae_minimum_links` and `ether802_3ad` arguments
* **data-source/junos_interface_physical**: remove same attributes in data source as resource
* **resource/junos_security_ipsec_vpn**: remove deprecated `bind_interface_auto` argument
* **resource/junos_system_radius_server**: remove deprecated `accouting_timeout` attribute

FEATURES:

* add **junos_forwardingoptions_sampling** resource (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))
* add **junos_oam_gretunnel_interface** resource (Fix [#457](https://github.com/jeremmfr/terraform-provider-junos/issues/457))
* add **junos_services_flowmonitoring_v9_template** resource (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))

ENHANCEMENTS:

* refactor provider to integrate new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
the resources and data sources will migrate progressively to this new plugin
* **provider**: add new cipher `aes256-gcm@openssh.com` in default ciphers when `ssh_ciphers` is not specified
* **provider**: enhance debug logs by prefixing messages with local and remote addresses of the sessions
* **resource/junos_forwardingoptions_sampling_instance**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `routing_instance` argument to allow create sampling instance in routing instance.  
  `id` attribute has now the format `<name>_-_<routing_instance>`
  * `flow_server` block argument is now unordered
* **resource/junos_services_flowmonitoring_vipfix_template**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `type` argument now accept `bridge-template`
  * add  `flow_key_output_interface`,  `mpls_template_label_position`, `template_refresh_rate` block argument (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))
* **resource/junos_interface_logical**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **data-source/junos_interface_logical**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource  
the data-source schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_interface_physical**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` arguments
* **data-source/junos_interface_physical**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource  
  the data-source schema has been upgraded to have one-blocks in single mode instead of list
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` attributes like resource
* **resource/junos_interface_physical_disable**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_interface_st0_unit**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **data-source/junos_interface_logical_info**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
the data-source schema has been upgraded to have one-blocks in single mode instead of list
* **data-source/junos_interfaces_physical_present**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value except `type` argument
* **data-source/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource
* **resource/junos_security**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `nat_source` block argument (Fix [#458](https://github.com/jeremmfr/terraform-provider-junos/issues/458))
* **resource/junos_security_address_book**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply (Fix [#403](https://github.com/jeremmfr/terraform-provider-junos/issues/403))  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_global_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_ike_gateway**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `distinguished_name_container` and `distinguished_name_wildcard` arguments inside `remote_identity` block argument
* **resource/junos_security_ike_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` and `reauth_frequency` arguments
* **resource/junos_security_ike_proposal**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_proposal**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_vpn**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `copy_outer_dscp`, `manual`, `multi_sa_forwarding_class` and `udp_encapsulate` arguments
* **resource/junos_security_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_policy_tunnel_pair_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_zone**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **data-source/junos_security_zone**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource
* **resource/junos_security_zone_book_address**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_zone_book_address_set**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional string attributes doesn't accept *empty* value

BUG FIXES:

* provider: when `ssh_retry_to_establish` > 1, stop retrying to open connections or sessions after a gracefully shutting down call with `Ctrl-c`
* **resource/junos_security_ike_gateway**: don't catch error when read `local_identity` and `remote_identity` block arguments and `type` is `distinguished-name`

## Previous Releases

* [v1](https://github.com/jeremmfr/terraform-provider-junos/blob/v1/CHANGELOG.md)
