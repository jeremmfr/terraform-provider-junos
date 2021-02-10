## upcoming release
ENHANCEMENTS:
* add `h323_disable`, `mgcp_disable`, `rtsp_disable`, `sccp_disable` and `sip_disable` arguments in `junos_security` resource (Fixes #95) Thanks [@a-d-v](https://github.com/a-d-v)
* add `default_address_selection` and `no_multicast_echo` arguments in `junos_system` resource (Fixes #97) Thanks [@a-d-v](https://github.com/a-d-v)
* add `junos_security_screen` resource (Fixes parts of [#92](https://github.com/jeremmfr/terraform-provider-junos/issues/92))
* add `junos_security_screen_whitelist` resource
* add `advance_policy_based_routing_profile`, `application_tracking`, `description`, `reverse_reroute`, `screen`, `source_identity_log` and `tcp_rst` arguments in `junos_security_zone` resource (Fixes parts of [#92](https://github.com/jeremmfr/terraform-provider-junos/issues/92))
* add `junos_security_utm_custom_url_category` resource (Fixes #108) Thanks [@a-d-v](https://github.com/a-d-v)

BUG FIXES:
* clean code: remove useless else when read a empty config
* fix typo in name of `accounting_timeout` argument in `junos_system_radius_server` resource. **Update your config for new version of this argument**
* fix warnings received from the device generate failures on resource actions. Now, received warnings are send to terraform under warnings format (Fixes #105)
* fix IP/Mask validation for point to point IPs

## 1.12.3 (February 5, 2021)
BUG FIXES:
* fix crash when `bind_interface` change in `junos_security_ipsec_vpn` resource

## 1.12.2 (February 3, 2021)
BUG FIXES:
* allow the name length of some objects > 32 for part of the resources (Fixes [#101](https://github.com/jeremmfr/terraform-provider-junos/issues/101))

## 1.12.1 (February 1, 2021)
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
* add `interface` option to `qualified_next_hop` on `static_route` resource (Fixes #71) Thanks [@tagur87](https://github.com/tagur87)
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
* fix `junos_interface` resource : crach on closeSession Netconf after error on startNewSession

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
