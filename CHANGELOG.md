## upcoming release
ENHANCEMENTS:
* add `alg` argument in resource `security` (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))
* add `flow` argument in resource `security` (Fixes parts of [#54](https://github.com/jeremmfr/terraform-provider-junos/issues/54))

BUG FIXES:

## v1.7.0
ENHANCEMENTS:
* add `dynamic_remote` argument in resource `security_ike_gateway` (Fixes [#50](https://github.com/jeremmfr/terraform-provider-junos/issues/50))
* add `aaa` argument in resource `security_ike_gateway`

BUG FIXES:
* fix lint errors from latest golangci-lint

## v1.6.1
BUG FIXES:
* fix compile libraries into release (for alpine linux like hashicorp/terraform docker image)

## v1.6.0
FEATURES:
* add resource `junos_security` (special resource for static configuration in security block) (Fixes [#43](https://github.com/jeremmfr/terraform-provider-junos/issues/43))
* add resource `junos_system` (special resource for static configuration in system block) (Fixes parts of [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add resource `junos_routing_options` (special resource for static configuration in routing-options block)

ENHANCEMENTS:
* add `sshkey_pem` argument in provider configuration
* add `send_mode` for `dead_peer_detection` in resource `security_ike_gateway` (Fixes [#43](https://github.com/jeremmfr/terraform-provider-junos/issues/43))
* upgrade to terraform-plugin-sdk v2
* switch to sdk for part of ValidateFunc and rewrite the others to ValidateDiagFunc
* code optimization (compact test err func if not nil)

BUG FIXES:
* fix sess.configLock return already nil

## v1.5.1
BUG FIXES:
* add missing `password` field in provider configuration for ssh authentication (Fixes [#41](https://github.com/jeremmfr/terraform-provider-junos/issues/41))

## v1.5.0
FEATURES:
* add data source `junos_interface`

ENHANCEMENTS:
* add argument `vlan_tagging_id` in resource junos_interface

## v1.4.0
FEATURES:
* add resource `junos_system_ntp_server` (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add resource `junos_system_radius_server` (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add resource `junos_system_syslog_host` (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))
* add resource `junos_system_syslog_file` (Fixes [#33](https://github.com/jeremmfr/terraform-provider-junos/issues/33))

ENHANCEMENTS:
* add options `apply_path`, `dynamic_db` in resource junos_policyoptions_prefix_list (Fixes [#31](https://github.com/jeremmfr/terraform-provider-junos/issues/31))
* add options `is_fragment`, `next_header`, `next_header_except` in `from` block for resource firewall_filter (Fixes [#32](https://github.com/jeremmfr/terraform-provider-junos/issues/32))

BUG FIXES:
* fix message validateIntRange

## v1.3.0
FEATURES:
* add resource `junos_security_utm_custom_url_pattern` (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add resource `junos_security_utm_policy` (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add resource `junos_security_utm_profile_web_filtering_juniper_enhanced` (Fixes [#26](https://github.com/jeremmfr/terraform-provider-junos/issues/26))
* add resource `junos_security_utm_profile_web_filtering_juniper_local`
* add resource `junos_security_utm_profile_web_filtering_websense_redirect`

ENHANCEMENTS:
* remove useless LF for list of set command

BUG FIXES:
* fix typo in errors and commits messages
* [workflows] fix compile freebsd/arm64 on release
* fix rule/policy with space in name for application-services in resource junos_security_policy
* fix no empty List if Required for many resource

## 1.2.1
ENHANCEMENTS:
for terraform 0.13
* upgrade go version
* [workflows] rewrite release job
* [doc] rewrite index/readme

BUG FIXES:
* [workflows] no tar.gz incompatible with registry

## 1.2.0
FEATURES:
* new resource `junos_aggregate_route` (Fixes [#24](https://github.com/jeremmfr/terraform-provider-junos/issues/24))

ENHANCEMENTS:
* add `community` on resource `junos_static_route`

BUG FIXES:
* fix go lint error

## 1.1.1
BUG FIXES:
* Allow usage of ~ in sshkeyfile path (Fixes [#22](https://github.com/jeremmfr/terraform-provider-junos/issues/22))

## 1.1.0
ENHANCEMENTS:
*  add `application-services` in `security_policy` (Fixes [#20](https://github.com/jeremmfr/terraform-provider-junos/issues/20))

## 1.0.6
BUG FIXES:
* update module go-netconf : Close ssh socket even if we get an error

## 1.0.5
BUG FIXES:
* fix resource `junos_interface` crach on closeSession Netconf after error on startNewSession

## 1.0.4
BUG FIXES:
* fix `bind_interface_auto` on resource `junos_security_ipsec_vpn` -> search st0 unit not in terse simply
* remove commit-check before commit which gives the same error if there is
* fix check interface disable and NC

## 1.0.3
BUG FIXES:
* fix terraform crash with an empty blocks-mode (no one required)

## 1.0.2
ENHANCEMENTS:
* move cmd/debug environnement variables to provider config

## 1.0.1
BUG FIXES:
* fix readInterface with empty/disappeared interface

## 1.0.0

First release
