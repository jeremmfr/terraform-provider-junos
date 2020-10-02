## upcoming release
ENHANCEMENTS:

BUG FIXES:
* add missing `password` field in provider configuration for ssh authentication #41

## v1.5.0
ENHANCEMENTS:
* add data source junos_interface
* add argument `vlan_tagging_id` in resource junos_interface

## v1.4.0
ENHANCEMENTS:
* add options `apply_path`, `dynamic_db` in resource junos_policyoptions_prefix_list #31
* add options `is_fragment`, `next_header`, `next_header_except` in `from` block for resource firewall_filter #32
* add resource junos_system_ntp_server #33
* add resource junos_system_radius_server #33
* add resource junos_system_syslog_host #33
* add resource junos_system_syslog_file #33

BUG FIXES:
* fix message validateIntRange

## v1.3.0
ENHANCEMENTS:
* add resource junos_security_utm_custom_url_pattern #26
* add resource junos_security_utm_policy #26
* add resource junos_security_utm_profile_web_filtering_juniper_enhanced #26
* add resource junos_security_utm_profile_web_filtering_juniper_local
* add resource junos_security_utm_profile_web_filtering_websense_redirect
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
ENHANCEMENTS:
* add community on resource junos_static_route
* new resource junos_aggregate_route #24

BUG FIXES:
* fix go lint error

## 1.1.1
BUG FIXES:
* Allow usage of ~ in sshkeyfile path (Fixes #22)

## 1.1.0
ENHANCEMENTS:
*  add application-services in security_policy (#20)

## 1.0.6
BUG FIXES:
* update module go-netconf : Close ssh socket even if we get an error

## 1.0.5
BUG FIXES:
* fix resource_junos_interface crach on closeSession Netconf after error on startNewSession  ([19](https://github.com/jeremmfr/terraform-provider-junos/pull/19))

## 1.0.4
BUG FIXES:
* fix ipsec_vpn bind_interface_auto -> search st0 unit not in terse simply ([17](https://github.com/jeremmfr/terraform-provider-junos/pull/17))
* remove commit-check before commit which gives the same error if there is ([16](https://github.com/jeremmfr/terraform-provider-junos/pull/16))
* fix check interface disable and NC ([15](https://github.com/jeremmfr/terraform-provider-junos/pull/15))

## 1.0.3
BUG FIXES:
* fix terraform crash with an empty blocks-mode (no one required) ([14](https://github.com/jeremmfr/terraform-provider-junos/pull/14))

## 1.0.2
ENHANCEMENTS:
* move cmd/debug environnement variables to provider config ([13](https://github.com/jeremmfr/terraform-provider-junos/pull/13))

## 1.0.1
BUG FIXES:
* fix readInterface with empty/disappeared interface

## 1.0.0

First release
