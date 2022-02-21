<!-- markdownlint-disable-file MD013 MD041 -->
## upcoming release

ENHANCEMENTS:

* resource/`junos_snmp`: add `engine_id` argument (Fixes parts of #339)
* add `junos_snmp_v3_usm_user` resource (Fixes parts of #339)
* add `junos_snmp_v3_vacm_accessgroup` resource (Fixes parts of #339)
* add `junos_snmp_v3_vacm_securitytogroup` resource (Fixes parts of #339)

BUG FIXES:

## 1.24.1 (February 11, 2022)

BUG FIXES:

* bump package github.com/jeremmfr/go-netconf to v0.4.2 to correctly detect protection on a private key in `sshkey_pem` or `sshkeyfile` and use passphrase from `keypass` argument (Fixes [#337](https://github.com/jeremmfr/terraform-provider-junos/issues/337))

## 1.24.0 (January 31, 2022)

FEATURES:

* add provider arguments `fake_update_also` and `fake_delete_also` to add the same workaround as with `fake_create_with_setfile`  
  **Don't use in normal terraform run** and **be careful with this option**  
  See docs for more information (Fixes [#329](https://github.com/jeremmfr/terraform-provider-junos/issues/329))

BUG FIXES:

* resource/`junos_security_policy`: read the config lines with `then allow tunnel pair-policy` to add them when updating the resource and thus avoid modifying the `junos_security_policy_tunnel_pair_policy` resource lines

## 1.23.0 (December 17, 2021)

ENHANCEMENTS:

* resource/`junos_security_address_book`: `network_address`, `wildcard_address`, `dns_name`, `range_address` and `address_set` block arguments are now unordered blocks. (Fixes [#316](https://github.com/jeremmfr/terraform-provider-junos/issues/316))
* resource/`junos_security_zone`: `address_book`, `address_book_dns`, `address_book_range`, `address_book_set` and `address_book_wildcard` block arguments are now unordered blocks. (Fixes [#169](https://github.com/jeremmfr/terraform-provider-junos/issues/169))

## 1.22.2 (December 07, 2021)

BUG FIXES:

* resource/`junos_interface`, `junos_interface_physical`, `junos_interface_logical`: fix detection of missing interfaces with the latest versions of Junos (the `not found` message has been moved from the reply field to the error field in the netconf message)
* resource/`junos_security_ike_gateway`: fix bad value used when adding Junos config with `wildcard` argument inside `distinguished_name` block inside `dynamic_remote` block
* resource/`junos_system_login_user`: fix validation of `name` with a dot that is acceptable for Junos (Fixes [#318](https://github.com/jeremmfr/terraform-provider-junos/issues/318))

## 1.22.1 (November 30, 2021)

BUG FIXES:

* resource/`junos_interface_logical`: add `srx_old_option_name` argument inside `dhcp` block argument to be able to fix the configuration of DHCP client on SRX devices with an old version of Junos (use option name `dhcp-client` instead of `dhcp` in `family inet`) (Fixes [#313](https://github.com/jeremmfr/terraform-provider-junos/issues/313))
* data-source/`junos_interface_logical`: add `srx_old_option_name` argument inside `dhcp` block argument (like resource)

## 1.22.0 (November 19, 2021)

FEATURES:

* add `junos_access_address_assignment_pool` resource (Fixes parts of [#301](https://github.com/jeremmfr/terraform-provider-junos/issues/301))
* add `junos_interface_physical_disable` resource (Fixes [#305](https://github.com/jeremmfr/terraform-provider-junos/issues/305))
* add `junos_interfaces_physical_present` data-source
* add `junos_system_services_dhcp_localserver_group` resource (Fixes parts of [#301](https://github.com/jeremmfr/terraform-provider-junos/issues/301))

ENHANCEMENTS:

* resource/`junos_application`: add `term` block argument (Fixes [#296](https://github.com/jeremmfr/terraform-provider-junos/issues/296)), add `inactivity_timeout_never` argument (Fixes [#308](https://github.com/jeremmfr/terraform-provider-junos/issues/308))
* resource/`junos_chassis_cluster`: add `control_ports` block argument (Fixes [#304](https://github.com/jeremmfr/terraform-provider-junos/issues/304))
* resource/`junos_group_dual_system`: add `system.0.inet6_backup_router_address` and `system.0.inet6_backup_router_destination` arguments, add validation on `system.0.backup_router_address` and list of string for `system.0.backup_router_destination` is now unordered (Fixes [#302](https://github.com/jeremmfr/terraform-provider-junos/issues/302))
* resource/`junos_interface_logical`: add `dhcp` and `dhcpv6_client` block arguments inside `family_inet` and `family_inet6` block arguments (Fixes parts of [#301](https://github.com/jeremmfr/terraform-provider-junos/issues/301))
* data-source/`junos_interface_logical`: add `dhcp` and `dhcpv6_client` block attributes inside `family_inet` and `family_inet6` block attributes
* resource/`junos_system`: add `ports` block argument (Fixes [#294](https://github.com/jeremmfr/terraform-provider-junos/issues/294))

BUG FIXES:

* resource/`junos_security_idp_custom_attack`: fix validation of IPv6 address for `destination_value`, `extension_header_destination_option_home_address_value` and `source_value` inside `protocol_ipv6` block
* resource/`junos_services_rpm_probe`: fix validation of IPv6 address for `inet6_source_address`, `rpm_scale.0.source_inet6_address_base`, `rpm_scale.0.source_inet6_step`, `rpm_scale.0.target_inet6_address_base` and `rpm_scale.0.target_inet6_step` inside `test` block

## 1.21.1 (October 22, 2021)

BUG FIXES:

* module go-netconf updated to enhance RPCError display with the `error-path` and `error-info>bad-element` values if set (Fixes parts of [#292](https://github.com/jeremmfr/terraform-provider-junos/issues/292))
* r/`*`: fix missing identifier value in the errors `multiple blocks with the same identifier` (Fixes parts of [#292](https://github.com/jeremmfr/terraform-provider-junos/issues/292))

## 1.21.0 (October 12, 2021)

FEATURES:

* add `junos_security_nat_static_rule` resource (Fixes [#281](https://github.com/jeremmfr/terraform-provider-junos/issues/281))

ENHANCEMENTS:

* resource/`*`: to avoid any confusion, the provider now detects and generates an error during `apply` when there are duplicate elements (with the same identifier, for example the same `name`) in the block lists of certain resources
* resource/`junos_routing_instance`: add `instance_export` and `instance_import` arguments (Fixes [#280](https://github.com/jeremmfr/terraform-provider-junos/issues/280))
* resource/`junos_routing_options`: add `instance_export` and `instance_import` arguments
* resource/`junos_security_address_book`: add `address_set` argument inside `address_set` block (Fixes [#287](https://github.com/jeremmfr/terraform-provider-junos/issues/287))
* resource/`junos_security_nat_destination`: add multiple arguments, `application`, `destination_...`, `protocol`, `source_...` in `rule` block and `description`
* resource/`junos_security_nat_destination_pool`: add `description` argument
* resource/`junos_security_nat_source`: add multiple arguments, `application`, `destination_...`, `protocol`, `source_...` in `match` block inside `rule` block and `description`
* resource/`junos_security_nat_source_pool`: add `description` argument
* resource/`junos_security_nat_static`: add multiple arguments, `destination_...`, `source_...` in `rule` block, `mapped_port...` in `then` block, `configure_rules_singly` and `description`. also add possibility to use `prefix-name` in `then.0.type`
* resource/`junos_security_zone`: add `address_set` argument inside `address_book_set` block
* resource/`junos_security_zone_book_address_set`: add `address_set` argument
* release now with golang 1.17 and replace the terraform sdk to a fork to avoid the note `Objects have changed outside of Terraform` with the empty string lists when create resources

BUG FIXES:

* resource/`junos_ospf_area`: fix missing set interface when `interface` block have only `name` set
* resource/`junos_security_nat_source`: fix panic when `match` block inside `rule` block is empty

## 1.20.0 (September 07, 2021)

FEATURES:

* add `junos_eventoptions_generate_event` resource (Fixes [#267](https://github.com/jeremmfr/terraform-provider-junos/issues/267))
* add `junos_security_dynamic_address_feed_server` resource (Fixes parts of [#273](https://github.com/jeremmfr/terraform-provider-junos/issues/273))
* add `junos_security_dynamic_address_name` resource (Fixes parts of [#273](https://github.com/jeremmfr/terraform-provider-junos/issues/273))

ENHANCEMENTS:

* resource/`junos_chassis_cluster`: add `preempt_delay`, `preempt_limit` and `preempt_period` arguments inside `redundancy_group` block list argument (Fixes [#270](https://github.com/jeremmfr/terraform-provider-junos/issues/270))
* resource/`junos_firewall_filter`: arguments with type list of string in block `term.*.from` are now unordered
* resource/`junos_interface_logical`: add `dad_disable` argument  inside `family_inet6` block argument (Fixes [#263](https://github.com/jeremmfr/terraform-provider-junos/issues/263))
* data-source/`junos_interface`, `junos_interface_logical`: `vrrp_group.*.authentication_key` is now a sensitive argument (like resource)
* data-source/`junos_interface_logical`: add `dad_disable` attributes as for the resource
* resource/`junos_interface_logical`: lists of string for `security_inbound_protocols` and `security_inbound_services` are now unordered
* resource/`junos_policyoptions_policy_statement`: arguments with type list of string (except `policy`) in block `term.*.from` and `term.*.to` are now unordered
* resource/`junos_security`: list of string for `ike_traceoptions.0.flag` is now unordered
* resource/`junos_security`: add validation on `name` argument inside `file` block inside `ike_traceoptions` block
* resource/`junos_security_global_policy`: arguments with type list of string in block `policy` are now unordered
* resource/`junos_security_nat_source`: arguments with type list of string in block `rule.*.match` are now unordered
* resource/`junos_security_policy`: arguments with type list of string in block `policy` are now unordered
* resource/`junos_security_screen`: lists of string for `tcp.0.syn_flood.0.whitelist.*.destination_address`, `tcp.0.syn_flood.0.whitelist.*.source_address` and `udp.0.flood.0.whitelist` are now unordered
* resource/`junos_security_screen_whitelist`: list of string for `address` is now unordered
* resource/`junos_security_zone`: lists of string for `inbound_protocols` and `inbound_services` are now unordered
* resource/`junos_system`: arguments with type list of string are now unordered (except `authentication_order`, `name_server` and `ssh.0.authentication_order`)
* resource/`junos_system`: add `ntp` block argument (Fixes [#261](https://github.com/jeremmfr/terraform-provider-junos/issues/261))
* resource/`junos_system`: add `netconf_traceoptions` block argument inside `services` block argument (Fixes [#262](https://github.com/jeremmfr/terraform-provider-junos/issues/262))
* resource/`*`: sets of string are now ordered before adding to Junos config to avoid unnecessary diff in commits
* docs: rewrite style for argument name and type
* docs: add attributes reference on resource

BUG FIXES:

* resource/`junos_security`: fix reading `size` argument inside `file` block inside `ike_traceoptions` block when number match a multiple of 1024 (example 1k, 1m, 1g)
* resource/`junos_security`: fix string format for `idp_security_package.0.automatic_start_time` to `YYYY-MM-DD.HH:MM:SS` to avoid unnecessary diff for Terraform when timezone of Junos device change
* resource/`junos_chassis_cluster`: fix possible crash in certain conditions when import this resource
* resource/`*`: add validation to some arguments which cannot contain a space character and thus avoid bugs when reading these arguments

## 1.19.0 (July 30, 2021)

FEATURES:

* add `junos_eventoptions_destination` resource
* add `junos_eventoptions_policy` resource (Fixes [#252](https://github.com/jeremmfr/terraform-provider-junos/issues/252))

ENHANCEMENTS:

* resource/`junos_application`: add `application_protocol`, `description`, `ether_type`, `rpc_program_number` and `uuid` arguments (Fixes [#255](https://github.com/jeremmfr/terraform-provider-junos/issues/255))

## 1.18.0 (July 09, 2021)

FEATURES:

* add `junos_bridge_domain` resource
* add `junos_evpn` resource (Fixes parts of [#131](https://github.com/jeremmfr/terraform-provider-junos/issues/131))
* add `junos_services_rpm_probe` resource (Fixes [#247](https://github.com/jeremmfr/terraform-provider-junos/issues/247))
* add `junos_switch_options` resource (Fixes parts of [#131](https://github.com/jeremmfr/terraform-provider-junos/issues/131))

ENHANCEMENTS:

* resource/`junos_bgp_group`: `authentication_key` is now a sensitive argument
* resource/`junos_bgp_neighbor`: `authentication_key` is now a sensitive argument
* resource/`junos_interface`, `junos_interface_logical`: `vrrp_group.*.authentication_key` is now a sensitive argument
* resource/`junos_policyoptions_policy_statement`: add `add_it_to_forwarding_table_export` argument (Fixes [#241](https://github.com/jeremmfr/terraform-provider-junos/issues/241))
* resource/`junos_routing_instance`: add `description`, `route_distinguisher`, `vrf_export`, `vrf_import`, `vrf_target`, `vrf_target_auto`, `vrf_target_export`, `vrf_target_import`, `vtep_source_interface`, `configure_rd_vrfopts_singly` and `configure_type_singly` arguments
* resource/`junos_routing_options`: add `forwarding_table_export_configure_singly` argument
* resource/`junos_routing_options`: add `router_id` argument
* resource/`junos_security_ike_gateway`: `aaa.0.client_password` is now a sensitive argument
* resource/`junos_system`: `archival_configuration.0.archive_site.*.password` is now a sensitive argument

BUG FIXES:

* resource/`junos_ospf`: fix missing mutex unlocking when read resource and checking routing-instance existence
* resource/`junos_security_nat_destination`: fix order issue on `from.0.value` list
* resource/`junos_security_nat_source`: fix order issue on `from.0.value` and `to.0.value` lists (Fixes [#243](https://github.com/jeremmfr/terraform-provider-junos/issues/243))
* resource/`junos_security_nat_static`: fix order issue on `from.0.value` list
* resource/`junos_system`: unescape the html entities for `announcement` argument inside `login` argument (Fixes parts of [#251](https://github.com/jeremmfr/terraform-provider-junos/issues/251))
* resource/`junos_system`: remove the potential double quotes for `ciphers` argument inside `services.0.ssh` argument (Fixes parts of [#251](https://github.com/jeremmfr/terraform-provider-junos/issues/251))
* resource/`junos_vlan`: fix order issue on `community_vlans` and `vlan_id_list` lists

## 1.17.0 (June 18, 2021)

FEATURES:

* add `junos_ospf` resource
* add `junos_security_idp_custom_attack` resource (Fixes parts of [#225](https://github.com/jeremmfr/terraform-provider-junos/issues/225))
* add `junos_security_idp_custom_attack_group` resource
* add `junos_security_idp_policy` resource (Fixes parts of [#225](https://github.com/jeremmfr/terraform-provider-junos/issues/225))

ENHANCEMENTS:

* provider: try multiple SSH authentication methods (key + password)
* provider: add `ssh_ciphers` argument to configure ciphers used in SSH connection
* provider: add support of SSH agent to SSH authentication (Fixes [#212](https://github.com/jeremmfr/terraform-provider-junos/issues/212))
* resource/`junos_application`: add `inactivity_timeout` argument (Fixes [#230](https://github.com/jeremmfr/terraform-provider-junos/issues/230))
* resource/`junos_group_dual_system`: add `preferred` and `primary` arguments inside `family_inet_address` and `family_inet6_address` arguments inside `interface_fxp0` argument (Fixes [#211](https://github.com/jeremmfr/terraform-provider-junos/issues/211))
* resource/`junos_interface_logical`: add `preferred` and `primary` arguments inside `address` argument inside `family_inet` and `family_inet6` arguments, add `vlan_no_compute` argument to disable the automatic compute of `vlan_id`
* data-source/`junos_interface_logical`: add `preferred` and `primary` attributes as for the resource
* resource/`junos_routing_options`, `junos_security`, `junos_services`, `junos_snmp`: add `clean_on_destroy` argument to clean static configuration when destroy the resource (Fixes [#227](https://github.com/jeremmfr/terraform-provider-junos/issues/227))
* resource/`junos_routing_options`: add `forwarding_table` argument (Fixes [#221](https://github.com/jeremmfr/terraform-provider-junos/issues/221))
* resource/`junos_security`: add `idp_security_package` and `idp_sensor_configuration` arguments (Fixes parts of [#225](https://github.com/jeremmfr/terraform-provider-junos/issues/225)), add `user_identification_auth_source` argument (Fixes [#238](https://github.com/jeremmfr/terraform-provider-junos/issues/238))
* resource/`junos_security_global_policy`: add `idp_policy` argument inside `permit_application_services` argument inside `policy` argument
* resource/`junos_security_policy`: add `idp_policy` argument inside `permit_application_services` argument inside `policy` argument
* resource/`junos_services_flowmonitoring_vipfix_template`: add `ip_template_export_extension` argument (Fixes [#229](https://github.com/jeremmfr/terraform-provider-junos/issues/229))
* resource/`junos_system`: add `radius_options_attributes_nas_ipaddress`, `radius_options_enhanced_accounting` and `radius_options_password_protocol_mschapv2` arguments (Fixes [#210](https://github.com/jeremmfr/terraform-provider-junos/issues/210)), add `archival_configuration` argument (Fixes [#231](https://github.com/jeremmfr/terraform-provider-junos/issues/231))

## 1.16.2 (May 28, 2021)

BUG FIXES:

* provider: fix XML error on commit with RPC reply without `<commit-results>` but different from `<ok/>` (Fixes [#223](https://github.com/jeremmfr/terraform-provider-junos/issues/223))
* resource/`junos_interface_logical`: disable set vlan-id on 'vlan.*' interface (Fixes parts of [#222](https://github.com/jeremmfr/terraform-provider-junos/issues/222))
* resource/`junos_vlan`: allow 'vlan.*' interface in `l3_interface` argument (Fixes parts of [#222](https://github.com/jeremmfr/terraform-provider-junos/issues/222))

## 1.16.1 (May 26, 2021)

BUG FIXES:

* resource/`junos_interface_logical`: disable set vlan-id on 'irb.*' interface (Fixes [#217](https://github.com/jeremmfr/terraform-provider-junos/issues/217))

## 1.16.0 (May 17, 2021)

FEATURES:

* add `junos_security_zone_book_address` resource (Fixes parts of [#192](https://github.com/jeremmfr/terraform-provider-junos/issues/192))
* add `junos_security_zone_book_address_set` resource (Fixes parts of [#192](https://github.com/jeremmfr/terraform-provider-junos/issues/192))
* add `junos_services_advanced_anti_malware_policy` resource
* add `junos_services_ssl_initiation_profile` resource
* add `junos_services_user_identification_ad_access_domain` resource
* add `junos_services_user_identification_device_identity_profile` resource (Fixes parts of [#189](https://github.com/jeremmfr/terraform-provider-junos/issues/189))

ENHANCEMENTS:

* resource/`junos_security`: add `policies` argument with `policy_rematch` and `policy_rematch_extensive` arguments inside (Fixes [#185](https://github.com/jeremmfr/terraform-provider-junos/issues/185)) Thanks  [@Sha-San-P](https://github.com/Sha-San-P)
* resource/`junos_security_address_book`: list of string for `address` argument inside `address_set` argument is now unordered
* resource/`junos_security_global_policy`: add `advanced_anti_malware_policy` argument inside `permit_application_services` argument
* resource/`junos_security_global_policy`: add `match_source_end_user_profile` argument inside `policy` argument
* resource/`junos_security_nat_source_pool`: add `address_pooling` argument (Fixes [#193](https://github.com/jeremmfr/terraform-provider-junos/issues/193)) Thanks [@edpio19](https://github.com/edpio19)
* resource/`junos_security_policy`: add `match_source_end_user_profile` argument inside `policy` argument
* resource/`junos_security_policy`: add `advanced_anti_malware_policy` argument inside `permit_application_services` argument
* resource/`junos_security_zone`: add `address_book_configure_singly` argument to disable management of address-book in this resource. (Fixes parts of [#192](https://github.com/jeremmfr/terraform-provider-junos/issues/192))
* resource/`junos_security_zone`: add `address_book_dns`, `address_book_range` and `address_book_wildcard` arguments and add `description` on existing `address_book_*` arguments
* resource/`junos_security_zone`: list of string for `address` argument inside `address_book_set` argument is now unordered
* resource/`junos_services`: add `advanced_anti_malware` argument (Fixes [#201](https://github.com/jeremmfr/terraform-provider-junos/issues/201))
* resource/`junos_services`: add `user_identification` argument (Fixes parts of [#189](https://github.com/jeremmfr/terraform-provider-junos/issues/189))
* resource/`junos_services`: `url_parameter` is now a sensitive argument
* resource/`junos_services`: `authentication_token`, `authentication_tls_profile` and `url` are now attributes (information read from Junos config) when not set in Terraform config. (Fixes [#200](https://github.com/jeremmfr/terraform-provider-junos/issues/200))
* resource/`junos_system`: add `web_management_http` and `web_management_https` arguments (Fixes [#173](https://github.com/jeremmfr/terraform-provider-junos/issues/173)) Thanks [@MerryPlant](https://github.com/MerryPlant)
* resource/`junos_system`: add `license` argument (Fixes [#205](https://github.com/jeremmfr/terraform-provider-junos/issues/205)) Thanks [@MerryPlant](https://github.com/MerryPlant)
* clean code: remove override of the lists of 1 map to handle directly the map
* clean code: fix lll linter errors with a var to map
* clean code: remove a useless override of map with himself and a useless option to compare different type

BUG FIXES:

* fix errors not generated with certain nested blocks empty and default integer argument = -1 in these blocks
* resource/`junos_services`: fix set/read/delete empty `application_identification` block to enable 'application-identification'

## 1.15.1 (April 23, 2021)

BUG FIXES:

* resource/`junos_security_global_policy`: fix `match_application` argument not required if `match_dynamic_application` is set and Junos version is > 19.1R1 (Fixes [#188](https://github.com/jeremmfr/terraform-provider-junos/issues/188))
* resource/`junos_security_policy`: fix `match_application` argument not required if `match_dynamic_application` is set and Junos version is > 19.1R1 (Fixes [#188](https://github.com/jeremmfr/terraform-provider-junos/issues/188))

## 1.15.0 (April 20, 2021)

FEATURES:

* add `junos_forwardingoptions_sampling_instance` resource (Fixes parts of [#165](https://github.com/jeremmfr/terraform-provider-junos/issues/165))
* add `junos_generate_route` resource
* add `junos_services` resource (Fixes parts of [#145](https://github.com/jeremmfr/terraform-provider-junos/issues/145), [#158](https://github.com/jeremmfr/terraform-provider-junos/issues/158))
* add `junos_services_flowmonitoring_vipfix_template` resource (Fixes parts of [#165](https://github.com/jeremmfr/terraform-provider-junos/issues/165))
* add `junos_services_proxy_profile` resource
* add `junos_services_security_intelligence_policy` resource (Fixes parts of [#145](https://github.com/jeremmfr/terraform-provider-junos/issues/145))
* add `junos_services_security_intelligence_profile` resource (Fixes parts of [#145](https://github.com/jeremmfr/terraform-provider-junos/issues/145))
* add `junos_snmp` resource (Fixes parts of [#170](https://github.com/jeremmfr/terraform-provider-junos/issues/170))
* add `junos_snmp_clientlist` resource
* add `junos_snmp_community` resource (Fixes parts of [#170](https://github.com/jeremmfr/terraform-provider-junos/issues/170))
* add `junos_snmp_view` resource

ENHANCEMENTS:

* resource/`junos_aggregate_route`: add `as_path_*` arguments, add support IPv6 Routes and simplify delete lines when update
* resource/`junos_bgp_group`: add `keep_all` and `keep_none` arguments
* resource/`junos_bgp_neighbor`: add `keep_all` and `keep_none` arguments
* resource/`junos_group_dual_system`: add `family_inet6_address` argument inside `interface_fxp0` argument (Fixes [#177](https://github.com/jeremmfr/terraform-provider-junos/issues/177))
* resource/`junos_interface_logical`: add `sampling_input` and `sampling_output` arguments in `family_inet` and `family_inet6` arguments (Fixes parts of [#165](https://github.com/jeremmfr/terraform-provider-junos/issues/165))
* data-source/`junos_interface_logical`: add `sampling_input` and `sampling_output` attributes in `family_inet` and `family_inet6` attributes
* resource/`junos_security`: add `forwarding_process` argument (Fixes parts of [#158](https://github.com/jeremmfr/terraform-provider-junos/issues/158))
* resource/`junos_security`: add `rsh_disable` and `sql_disable` arguments (Fixes [#182](https://github.com/jeremmfr/terraform-provider-junos/issues/182)) Thanks [@edpio19](https://github.com/edpio19)
* resource/`junos_security_global_policy`: add `match_destination_address_excluded` and `match_source_address_excluded` arguments (Fixes [#159](https://github.com/jeremmfr/terraform-provider-junos/issues/159))
* resource/`junos_security_global_policy`: add `match_dynamic_application` arguments (Fixes parts of [#158](https://github.com/jeremmfr/terraform-provider-junos/issues/158))
* resource/`junos_security_nat_source_pool`: add `pool_utilization_alarm_raise_threshold` and `pool_utilization_alarm_clear_threshold` arguments (Fixes [#171](https://github.com/jeremmfr/terraform-provider-junos/issues/171)) Thanks [@edpio19](https://github.com/edpio19)
* resource/`junos_security_policy`: add `match_destination_address_excluded` and `match_source_address_excluded` arguments (Fixes [#159](https://github.com/jeremmfr/terraform-provider-junos/issues/159))
* resource/`junos_security_policy`: add `match_dynamic_application` arguments (Fixes parts of [#158](https://github.com/jeremmfr/terraform-provider-junos/issues/158))
* resource/`junos_static_route`: add `as_path_*` arguments and simplify delete lines when update
* resource/`junos_vlan`: add `vni_extend_evpn` argument inside `vxlan` argument (Fixes [#132](https://github.com/jeremmfr/terraform-provider-junos/issues/132)) Thanks [@dejongm](https://github.com/dejongm)
* clean code: remove useless type/func exporting, and fixes formatting golang code

BUG FIXES:

* fix panic when candidate config clear or unlock generate Junos error(s)
* fix missing MinItems=1 on a part of required `ListOfString` arguments

## 1.14.0 (March 19, 2021)

FEATURES:

* add `junos_chassis_cluster` resource (Fixes parts of [#106](https://github.com/jeremmfr/terraform-provider-junos/issues/106))
* add `junos_group_dual_system` resource (Fixes [#120](https://github.com/jeremmfr/terraform-provider-junos/issues/120))
* add `junos_null_commit_file` resource (Fixes parts of [#136](https://github.com/jeremmfr/terraform-provider-junos/issues/136))
* add `junos_security_address_book` resource (Fixes [#137](https://github.com/jeremmfr/terraform-provider-junos/issues/137)) Thanks [@tagur87](https://github.com/tagur87)
* add `junos_security_global_policy` resource (Fixes [#138](https://github.com/jeremmfr/terraform-provider-junos/issues/138))
* add provider argument `file_permission`
* add provider argument `fake_create_with_setfile` -  **Don't use in normal terraform run** and **be carefully with this option**
See docs for more information (Fixes parts of [#136](https://github.com/jeremmfr/terraform-provider-junos/issues/136))

ENHANCEMENTS:

* add `cluster`, `family_evpn` arguments in `junos_bgp_group` and `junos_bgp_neighbor` resource
* add new `bgp_multipath` block argument to replace `multipath` bool argument in `junos_bgp_group` and `junos_bgp_neighbor` resource
`bgp_multipath` let add optional arguments. `multipath` is now **deprecated**
* add `esi` argument in `junos_interface_physical` resource and data source (Fixes [#126](https://github.com/jeremmfr/terraform-provider-junos/issues/126)) Thanks [@dejongm](https://github.com/dejongm)
* add `ether_opts`, `gigether_opts` and `parent_ether_opts` arguments in `junos_interface_physical` resource and data source to add more options and replace `ae_lacp`, `ae_link_speed`, `ae_minimum_links`, `ether802_3ad` arguments which are now deprecated (Fixes [#133](https://github.com/jeremmfr/terraform-provider-junos/issues/133), [#127](https://github.com/jeremmfr/terraform-provider-junos/issues/127), parts of [#106](https://github.com/jeremmfr/terraform-provider-junos/issues/106))
* add `security_inbound_protocols` and `security_inbound_services` arguments in `junos_interface_logical` resource and data source (Fixes [#141](https://github.com/jeremmfr/terraform-provider-junos/issues/141))
* add `feature_profile_web_filtering_juniper_enhanced_server` argument in `utm` argument of `junos_security` resource (Fixes [#155](https://github.com/jeremmfr/terraform-provider-junos/issues/155))

BUG FIXES:

* fix change `description` to null in `junos_interface_logical` and `junos_interface_physical` resource
* fix `prefix` list order issue in `junos_policyoptions_prefix_list` resource (Fixes [#150](https://github.com/jeremmfr/terraform-provider-junos/issues/150))
* fix validation for `name` of `address_book` and `address_boob_set` in `junos_security_zone` resource (Fixes [#153](https://github.com/jeremmfr/terraform-provider-junos/issues/153))

## 1.13.1 (February 18, 2021)

BUG FIXES:

* fix source nat pool network address not allowed (Fixes [#128](https://github.com/jeremmfr/terraform-provider-junos/issues/128))

## 1.13.0 (February 11, 2021)

FEATURES:

* add `junos_security_screen` resource (Fixes parts of [#92](https://github.com/jeremmfr/terraform-provider-junos/issues/92))
* add `junos_security_screen_whitelist` resource
* add `junos_security_utm_custom_url_category` resource (Fixes [#108](https://github.com/jeremmfr/terraform-provider-junos/issues/108)) Thanks [@a-d-v](https://github.com/a-d-v)

ENHANCEMENTS:

* add `h323_disable`, `mgcp_disable`, `rtsp_disable`, `sccp_disable` and `sip_disable` arguments in `junos_security` resource (Fixes [#95](https://github.com/jeremmfr/terraform-provider-junos/issues/95)) Thanks [@a-d-v](https://github.com/a-d-v)
* add `default_address_selection` and `no_multicast_echo` arguments in `junos_system` resource (Fixes [#97](https://github.com/jeremmfr/terraform-provider-junos/issues/97)) Thanks [@a-d-v](https://github.com/a-d-v)
* add `advance_policy_based_routing_profile`, `application_tracking`, `description`, `reverse_reroute`, `screen`, `source_identity_log` and `tcp_rst` arguments in `junos_security_zone` resource (Fixes parts of [#92](https://github.com/jeremmfr/terraform-provider-junos/issues/92))

BUG FIXES:

* fix typo in name of `accounting_timeout` argument in `junos_system_radius_server` resource. **Update your config for new version of this argument**
* fix warnings received from the device generate failures on resource actions. Now, received warnings are send to terraform under warnings format (Fixes [#105](https://github.com/jeremmfr/terraform-provider-junos/issues/105))
* fix possibility to create `junos_interface_physical` and `junos_interface_logical` resource on a non-existent interface (Fixes [#111](https://github.com/jeremmfr/terraform-provider-junos/issues/111)). Read configuration before read interface status for validate resource existence.
* fix integer compute for `chassis aggregated-devices ethernet device-count` when create/update/delete `junos_interface_physical` resource. Now this uses current configuration instead of the status of 'ae' interfaces and also takes into account resource with prefix name 'ae' in addition to `ether802_3ad` argument.
* fix `filter_output` not set with good argument for `family inet6` in `junos_interface_logical` resource (Fixes [#117](https://github.com/jeremmfr/terraform-provider-junos/issues/117))
* fix IP/Mask validation for point to point IPs
* clean code: remove useless else when read a empty config

## 1.12.3 (February 05, 2021)

BUG FIXES:

* fix crash when `bind_interface` change in `junos_security_ipsec_vpn` resource

## 1.12.2 (February 03, 2021)

BUG FIXES:

* allow the name length of some objects > 32 for part of the resources (Fixes [#101](https://github.com/jeremmfr/terraform-provider-junos/issues/101))

## 1.12.1 (February 01, 2021)

BUG FIXES:

* possible mismatch for routing_instance in junos_interface_logical resource (Fixes [#98](https://github.com/jeremmfr/terraform-provider-junos/issues/98))
* can't create empty junos_policyoptions_prefix_list resource (Fixes [#99](https://github.com/jeremmfr/terraform-provider-junos/issues/99))

## 1.12.0 (January 20, 2021)

FEATURES:

* add `junos_system_login_class` resource (Fixes parts of [#88](https://github.com/jeremmfr/terraform-provider-junos/issues/88))
* add `junos_system_login_user` resource (Fixes parts of [#88](https://github.com/jeremmfr/terraform-provider-junos/issues/88))
* add `junos_system_root_authentication` resource

ENHANCEMENTS:

* add `ssh_sleep_closed` argument in provider configuration (Fixes part of [#87](https://github.com/jeremmfr/terraform-provider-junos/issues/87))
* add `login` argument in `junos_system` resource (Fixes parts of [#88](https://github.com/jeremmfr/terraform-provider-junos/issues/88))

BUG FIXES:

* add missing lock in data source to reduce netconf commands parallelism
* use only one ssh connection per action and per resource (Fixes part of [#87](https://github.com/jeremmfr/terraform-provider-junos/issues/87))

## 1.11.0 (January 05, 2021)

FEATURES:

* add `junos_interface_physical` resource for replace the parts of physical interface in deprecated `junos_interface` resource
* add `junos_interface_physical` data source for replace the parts of physical interface in deprecated `junos_interface` data source
* add `junos_interface_logical` resource for replace the parts of logical interface in deprecated `junos_interface` resource
* add `junos_interface_logical` data source for replace the parts of logical interface in deprecated `junos_interface` data source

ENHANCEMENTS:

* add `authentication_order`, `auto_snapshot`, `domain_name`, `host_name`, `inet6_backup_router`, `internet_options`, `max_configuration_rollbacks`, `max_configurations_on_flash`, `no_ping_record_route`, `no_ping_time_stamp`, `no_redirects`, `no_redirects_ipv6` and `time_zone` arguments in `junos_system` resource (Fixes [#81](https://github.com/jeremmfr/terraform-provider-junos/issues/81))
* code optimization (remove useless list length check before loop on)
* code optimization (remove useless strings mod usage to compare fixed string)
* deprecate `junos_interface` resource for two new resources (split physical and logical interface into separate resources)
* deprecate `junos_interface` data source for two new data sources (split physical and logical interface into separate data sources)

BUG FIXES:

* generate errors on apply if `syslog`, `services` or `services.0.ssh` block is set but empty in `junos_system` resource

## 1.10.0 (December 15, 2020)

ENHANCEMENTS:

* add `interface` option to `qualified_next_hop` on `static_route` resource (Fixes [#71](https://github.com/jeremmfr/terraform-provider-junos/issues/71)) Thanks [@tagur87](https://github.com/tagur87)
* add `inet_rpf_check` and `inet6_rpf_check` arguments in `junos_interface` resource (Fixes [#72](https://github.com/jeremmfr/terraform-provider-junos/issues/72))
* add `discard`, `receive`, `reject`, `next_table`, `active`, `passive`, `install`, `no_install`, `readvertise`, `no_readvertise`, `resolve`, `no_resolve`, `retain` and `no_retain` arguments in `junos_static_route` resource

BUG FIXES:

* fix missing compatibility argument checks when apply `junos_interface` resource (unit interface or not)
* fix `advertisements_threshold` argument missing for vrrp in family inet6 address in `junos_interface` resource

## 1.9.0 (December 03, 2020)

FEATURES:

* add `junos_system_information` data source (Fixes [#60](https://github.com/jeremmfr/terraform-provider-junos/issues/60)) Thanks [@tagur87](https://github.com/tagur87)
* add `junos_interface_st0_unit` resource (Fixes [#64](https://github.com/jeremmfr/terraform-provider-junos/issues/64))

ENHANCEMENTS:

* simplify gather system/software information when create new netconf session
* add support static IPv6 Routes in `junos_static_route` resource (Fixes [#67](https://github.com/jeremmfr/terraform-provider-junos/issues/67))

BUG FIXES:

* fix inconsistent result after creating `junos_interface` resource with only `name` argument (Fixes [#65](https://github.com/jeremmfr/terraform-provider-junos/issues/65))

## 1.8.0 (November 20, 2020)

FEATURES:

* add `junos_security_log_stream` resource (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))

ENHANCEMENTS:

* add `traffic_selector` argument in `junos_security_ipsec_vpn` resource (Fixes [#53](https://github.com/jeremmfr/terraform-provider-junos/issues/53))
* add `complete_destroy` argument in `junos_interface` resource
* add `alg` argument in `junos_security` resource (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))
* add `flow` argument in `junos_security` resource (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))
* add `log` argument in `junos_security` resource (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))
* add `forwarding_options` argument in `junos_security` resource (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))
* add `proposal_set` argument in `junos_security_ike_policy` and `junos_security_ipsec_policy` resource (Fixes [#55](https://github.com/jeremmfr/terraform-provider-junos/issues/55))
* add `icmp_code` and `icmp_code_except` sub-arguments for 'term.N.from' to `junos_firewall_filter` resource (Fixes [#58](https://github.com/jeremmfr/terraform-provider-junos/issues/58))
* optimize memory usage of functions for bgp_* resource
* release now with golang 1.15

BUG FIXES:

* remove useless ForceNew for `bind_interface_auto` argument in `junos_security_ipsec_vpn` resource

## 1.7.0 (November 03, 2020)

ENHANCEMENTS:

* add `dynamic_remote` argument in `junos_security_ike_gateway` resource (Fixes [#50](https://github.com/jeremmfr/terraform-provider-junos/issues/50))
* add `aaa` argument in `junos_security_ike_gateway` resource

BUG FIXES:

* fix lint errors from latest golangci-lint

## 1.6.1 (October 22, 2020)

BUG FIXES:

* fix compile libraries into release (for alpine linux like hashicorp/terraform docker image)

## 1.6.0 (October 21, 2020)

FEATURES:

* add `junos_security` resource (special resource for static configuration in security block) (Fixes [#43](https://github.com/jeremmfr/terraform-provider-junos/issues/43))
* add `junos_system` resource (special resource for static configuration in system block) (Fixes parts of [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add `junos_routing_options` resource (special resource for static configuration in routing-options block)

ENHANCEMENTS:

* add `sshkey_pem` argument in provider configuration
* add `send_mode` for `dead_peer_detection` in `junos_security_ike_gateway` resource (Fixes [#43](https://github.com/jeremmfr/terraform-provider-junos/issues/43))
* upgrade to terraform-plugin-sdk v2
* switch to sdk for part of ValidateFunc and rewrite the others to ValidateDiagFunc
* code optimization (compact test err func if not nil)

BUG FIXES:

* fix sess.configLock return already nil

## 1.5.1 (October 02, 2020)

BUG FIXES:

* add missing `password` field in provider configuration for ssh authentication (Fixes [#41](https://github.com/jeremmfr/terraform-provider-junos/issues/41))

## 1.5.0 (September 14, 2020)

FEATURES:

* add `junos_interface` data source

ENHANCEMENTS:

* add `vlan_tagging_id` argument in `junos_interface` resource

## 1.4.0 (September 04, 2020)

FEATURES:

* add `junos_system_ntp_server` resource (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add `junos_system_radius_server` resource (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add `junos_system_syslog_host` resource (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add `junos_system_syslog_file` resource (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))

ENHANCEMENTS:

* add `apply_path`, `dynamic_db` arguments in `junos_policyoptions_prefix_list` resource (Fixes [#31](https://github.com/jeremmfr/terraform-provider-junos/issues/31))
* add `is_fragment`, `next_header`, `next_header_except` arguments in `from` block for `junos_firewall_filter` resource (Fixes [#32](https://github.com/jeremmfr/terraform-provider-junos/issues/32))

BUG FIXES:

* fix message validateIntRange

## 1.3.0 (August 24, 2020)

FEATURES:

* add `junos_security_utm_custom_url_pattern` resource (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add `junos_security_utm_policy` resource (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add `junos_security_utm_profile_web_filtering_juniper_enhanced` resource (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add `junos_security_utm_profile_web_filtering_juniper_local` resource
* add `junos_security_utm_profile_web_filtering_websense_redirect` resource

ENHANCEMENTS:

* remove useless LF for list of set command

BUG FIXES:

* fix typo in errors and commits messages
* [workflows] fix compile freebsd/arm64 on release
* fix rule/policy with space in name for application-services in `junos_security_policy` resource
* fix no empty List if Required for many resource

## 1.2.1 (August 17, 2020)

ENHANCEMENTS:

for terraform 0.13

* upgrade go version
* [workflows] rewrite release job
* [doc] rewrite index/readme

BUG FIXES:

* [workflows] no tar.gz incompatible with registry

## 1.2.0 (July 21, 2020)

FEATURES:

* add `junos_aggregate_route` resource (Fixes [#24](https://github.com/jeremmfr/terraform-provider-junos/issues/24))

ENHANCEMENTS:

* add `community` argument on `junos_static_route` resource

BUG FIXES:

* fix go lint error

## 1.1.1 (June 28, 2020)

BUG FIXES:

* allow usage of ~ in sshkeyfile path (Fixes [#22](https://github.com/jeremmfr/terraform-provider-junos/issues/22))

## 1.1.0 (June 17, 2020)

ENHANCEMENTS:

* add `application-services` argument in `junos_security_policy` resource (Fixes [#20](https://github.com/jeremmfr/terraform-provider-junos/issues/20))

## 1.0.6 (May 28, 2020)

BUG FIXES:

* update module go-netconf : Close ssh socket even if we get an error

## 1.0.5 (March 26, 2020)

BUG FIXES:

* fix `junos_interface` resource : crash on closeSession Netconf after error on startNewSession

## 1.0.4 (January 03, 2020)

BUG FIXES:

* fix `bind_interface_auto` argument on `junos_security_ipsec_vpn` resource -> search st0 unit not in terse simply
* remove commit-check before commit which gives the same error if there is
* fix check interface disable and NC

## 1.0.3 (January 03, 2020)

BUG FIXES:

* fix terraform crash with an empty blocks-mode (no one required)

## 1.0.2 (January 03, 2020)

ENHANCEMENTS:

* move cmd/debug environnement variables to provider config

## 1.0.1 (December 18, 2019)

BUG FIXES:

* fix readInterface with empty/disappeared interface

## 1.0.0 (November 27, 2019)

First release
