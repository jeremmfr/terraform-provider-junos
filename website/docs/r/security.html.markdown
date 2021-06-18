---
layout: "junos"
page_title: "Junos: junos_security"
sidebar_current: "docs-junos-resource-security"
description: |-
  Configure static configuration in security block (when Junos device supports it)
---

# junos_security

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `security` block. By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `security` block

## Example Usage

```hcl
# Configure security
resource junos_security "security" {
  ike_traceoptions {
    file {
      name  = "ike.log"
      files = 5
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `clean_on_destroy` - (Optional)(`Bool`) Clean supported lines when destroy this resource.
* `alg` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'alg' configuration. See the [`alg` arguments] (#alg-arguments) block.
* `flow` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow' configuration. See the [`flow` arguments] (#flow-arguments) block.
* `forwarding_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'forwarding-options' configuration. See the [`forwarding_options` arguments] (#forwarding_options-arguments) block.
* `forwarding_process` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'forwarding-process' configuration.
  * `enhanced_services_mode` - (Optional)(`Bool`) Enable enhanced application services mode.
* `idp_security_package` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'idp security-package' configuration. See the [`idp_security_package` arguments] (#idp_security_package-arguments) block.
* `idp_sensor_configuration` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'idp sensor-configuration' configuration. See the [`idp_sensor_configuration` arguments] (#idp_sensor_configuration-arguments) block.
* `ike_traceoptions` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'ike traceoptions' configuration.
  * `file` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'file' configuration. See the [`file` arguments for ike_traceoptions] (#file-arguments-for-ike_traceoptions) block.
  * `flag` - (Optional)(`ListOfString`) Tracing parameters for IKE.
  * `rate_limit` - (Optional)(`Int`) Limit the incoming rate of trace messages (0..4294967295)
* `log` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'log' configuration. See the [`log` arguments] (#log-arguments) block.
* `policies` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'policies' configuration.
  * `policy_rematch` - (Optional)(`Bool`) Can be specified to allow session to remain open when an associated security policy is modified. Conflict with `policy_rematch_extensive`.
  * `policy_rematch_extensive` - (Optional)(`Bool`) Can be specified to allow session to remain open when an associated security policy is modified, renamed, deactivated, or deleted. Conflict with `policy_rematch`.
* `user_identification_auth_source` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'user-identification authentication-source' configuration.
  * `ad_auth_priority` - (Optional)(`Int`) Active-directory-authentication-table priority. Larger number means lower priority, 0 for disable (0..65535).
  * `aruba_clearpass_priority` - (Optional)(`Int`) ClearPass-authentication-table priority. Larger number means lower priority, 0 for disable (0..65535).
  * `firewall_auth_priority` - (Optional)(`Int`) Firewall-authentication priority. Larger number means lower priority, 0 for disable (0..65535).
  * `local_auth_priority` - (Optional)(`Int`) Local-authentication-table priority. Larger number means lower priority, 0 for disable (0..65535).
  * `unified_access_control_priority` - (Optional)(`Int`) Unified-access-control priority. Larger number means lower priority, 0 for disable (0..65535).
* `utm` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'utm' configuration.
  * `feature_profile_web_filtering_type` - (Optional)(`String`) Configuring feature-profile web-filtering type. Need to be 'juniper-enhanced', 'juniper-local', 'web-filtering-none' or 'websense-redirect'.
  * `feature_profile_web_filtering_juniper_enhanced_server` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'server' configuration. See the [`feature_profile_web_filtering_juniper_enhanced_server` arguments for utm] (#feature_profile_web_filtering_juniper_enhanced_server-arguments-for-utm) block.

---
#### alg arguments
* `dns_disable` - (Optional)(`Bool`) Disable dns alg.
* `ftp_disable` - (Optional)(`Bool`) Disable ftp alg.
* `h323_disable` - (Optional)(`Bool`) Disable h323 alg.
* `mgcp_disable` - (Optional)(`Bool`) Disable mgcp alg.
* `msrpc_disable` - (Optional)(`Bool`) Disable msrpc alg.
* `pptp_disable` - (Optional)(`Bool`) Disable pptp alg.
* `rsh_disable` - (Optional)(`Bool`) Disable rsh alg.
* `rtsp_disable` - (Optional)(`Bool`) Disable rtsp alg.
* `sccp_disable` - (Optional)(`Bool`) Disable sccp alg.
* `sip_disable` - (Optional)(`Bool`) Disable sip alg.
* `sql_disable` - (Optional)(`Bool`) Disable sql alg.
* `sunrpc_disable` - (Optional)(`Bool`) Disable sunrpc alg.
* `talk_disable` - (Optional)(`Bool`) Disable talk alg.
* `tftp_disable` - (Optional)(`Bool`) Disable tftp alg.

---
#### file arguments for ike_traceoptions
* `name` - (Optional)(`String`) Name of file in which to write trace information.
* `files` - (Optional)(`Int`) Maximum number of trace files (2..1000).
* `match` - (Optional)(`String`) Regular expression for lines to be logged.
* `no_world_readable` - (Optional)(`Bool`) Don't allow any user to read the log file.
* `size` - (Optional)(`Int`) Maximum trace file size (10240..1073741824)
* `world_readable` - (Optional)(`Bool`) Allow any user to read the log file

---
#### flow arguments
* `advanced_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow advanced-options' configuration.
  * `drop_matching_reserved_ip_address` - (Optional)(`Bool`) Drop matching reserved source IP address.
  * `drop_matching_link_local_address` - (Optional)(`Bool`) Drop matching link local address.
  * `reverse_route_packet_mode_vr` - (Optional)(`Bool`) Allow reverse route lookup with packet mode vr.
* `aging` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow aging' configuration.
  * `early_ageout` - (Optional)(`Int`) Delay before device declares session invalid (1..65535 seconds).
  * `high_watermark` - (Optional)(`Bool`) Percentage of session-table capacity at which aggressive aging-out starts (0..100 percent).
  * `low_watermark` - (Optional)(`Bool`) Percentage of session-table capacity at which aggressive aging-out ends (0..100 percent).
* `allow_dns_reply` - (Optional)(`Bool`) Allow unmatched incoming DNS reply packet.
* `allow_embedded_icmp` - (Optional)(`Bool`) Allow embedded ICMP packets not matching a session to pass through.
* `allow_reverse_ecmp` - (Optional)(`Bool`) Allow reverse ECMP route lookup.
* `enable_reroute_uniform_link_check_nat` - (Optional)(`Bool`) Enable reroute check with uniform link and NAT check.
* `ethernet_switching` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow ethernet-switching' configuration.
  * `block_non_ip_all` - (Optional)(`Bool`) Block all non-IP and non-ARP traffic including broadcast/multicast.
  * `bypass_non_ip_unicast` - (Optional)(`Bool`) Allow all non-IP (including unicast) traffic.
  * `bpdu_vlan_flooding` - (Optional)(`Bool`) Set 802.1D BPDU flooding based on VLAN.
  * `no_packet_flooding` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for stop IP flooding, send ARP/ICMP to trigger MAC learning.  
  There is one argument : `no_trace_route` - (Optional)(`Bool`) Don't send ICMP to trigger MAC learning.
* `force_ip_reassembly` - (Optional)(`Bool`) Force to reassemble ip fragments.
* `ipsec_performance_acceleration` - (Optional)(`Bool`) Accelerate the IPSec traffic performance.
* `mcast_buffer_enhance` - (Optional)(`Bool`) Allow to hold more packets during multicast session creation.
* `pending_sess_queue_length` - (Optional)(`String`) Maximum queued length per pending session. Need to be 'high', 'moderate' or 'normal'.
* `preserve_incoming_fragment_size` - (Optional)(`Bool`) Preserve incoming fragment size for egress MTU.
* `route_change_timeout` - (Optional)(`Int`) Timeout value for route change to nonexistent route (6..1800 seconds).
* `syn_flood_protection_mode` - (Optional)(`String`) TCP SYN flood protection mode. Need to be 'syn-cookie' or 'syn-proxy'.
* `sync_icmp_session` - (Optional)(`Bool`) Allow icmp sessions to sync to peer node.
* `tcp_mss` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow tcp-mss' configuration.
  * `all_tcp_mss` - (Optional)(`Int`) Enable MSS override for all packets with this value.
  * `gre_in` - Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for enable MSS override for all GRE packets coming out of an IPSec tunnel.
  There is one argument : `mss` - (Optional)(`Int`) MSS Value.
  * `gre_out` - Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for enable MSS override for all GRE packets entering an IPsec tunnel.
  There is one argument : `mss` - (Optional)(`Int`) MSS Value.
  * `ipsec_vpn` - Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for enable MSS override for all packets entering IPSec tunnel.
  There is one argument : `mss` - (Optional)(`Int`) MSS Value.
* `tcp_session` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'flow tcp-session' configuration.
  * `fin_invalidate_session` - (Optional)(`Bool`) Immediately end session on receipt of fin (FIN) segment.
  * `maximum_window` -  Maximum TCP proxy scaled receive window. Need to be '64K', '128K', '256K', '512K' or '1M'.
  * `no_sequence_check` - (Optional)(`Bool`) Disable sequence-number checking.
  * `no_syn_check` - (Optional)(`Bool`) Disable creation-time SYN-flag check. Conflict with `strict_syn_check`.
  * `no_syn_check_in_tunnel` - (Optional)(`Bool`) Disable creation-time SYN-flag check for tunnel packets. Conflict with `strict_syn_check`.
  * `rst_invalidate_session` - (Optional)(`Bool`) Immediately end session on receipt of reset (RST) segment.
  * `rst_sequence_check` - (Optional)(`Bool`) Check sequence number in reset (RST) segment.
  * `strict_syn_check` - (Optional)(`Bool`) Enable strict syn check. Conflict with `no_sync_check` and `no_syn_check_in_tunnel`.
  * `tcp_initial_timeout` - (Optional)(`Int`) Timeout for TCP session when initialization fails (4..300 seconds).
  * `time_wait_state` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare session timeout value in time-wait state. See the [`time_wait_state` arguments for tcp_session in flow] (#time_wait_state-arguments-for-tcp_session-in-flow) block.

---
#### forwarding_options arguments
* `inet6_mode` - (Optional)(`String`) Forwarding mode for inet6 family. Need to be 'drop', 'flow-based' or 'packet-based'.
* `mpls_mode` - (Optional)(`String`) Forwarding mode for mpls family. Need to be 'flow-based' or 'packet-based'.
* `iso_mode_packet_based` - (Optional)(`Bool`) Forwarding mode packet-based for iso family.

---
#### idp_security_package arguments
* `automatic_enable` - (Optional)(`Bool`) Enable scheduled download and update.
* `automatic_interval` - (Optional)(`Int`) Automatic interval (1..336 hours).
* `automatic_start_time` - (Optional)(`String`) Automatic start time (YYYY-MM-DD.HH:MM:SS +ZZZZ).
* `install_ignore_version_check` - (Optional)(`Bool`) Skip version check  when attack database gets installed.
* `proxy_profile` - (Optional)(`String`) Proxy profile of security package download.
* `source_address` - (Optional)(`String`) Source address to be used for sending download request.
* `url` - (Optional)(`String`) URL of Security package download.

---
#### idp_sensor_configuration arguments
* `log_cache_size` - (Optional)(`Int`) Log cache size (1..65535).
* `log_suppression` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'log suppression'.
  * `disable` - (Optional)(`Bool`) Disable log suppression.
  * `include_destination_address` - (Optional)(`Bool`) Include destination address while performing a log suppression.
  * `no_include_destination_address` - (Optional)(`Bool`) Don't include destination address while performing a log suppression.
  * `max_logs_operate` - (Optional)(`Int`) Maximum logs can be operate on (256..65536).
  * `max_time_report` - (Optional)(`Int`) Time after suppressed logs will be reported (1..60).
  * `start_log` - (Optional)(`Int`) Suppression start log (1..128).
* `packet_log` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'packet-log' configuration.
  * `source_address` - (Required)(`String`) Source IP address used to transport packetlog to a host.
  * `host_address` - (Optional)(`String`) Destination host to send packetlog to.
  * `host_port` - (Optional)(`Int`) Destination UDP port number (1..65536).
  * `max_sessions` - (Optional)(`Int`) Max num of sessions in unit(%) (1..100).
  * `threshold_logging_interval` - (Optional)(`Int`) Interval of logs for max limit session/memory reached in minutes (1..60).
  * `total_memory` - (Optional)(`Int`) Total memory unit(%) (1..100).
* `security_configuration_protection_mode` - (Optional)(`String`) Enable security protection mode.

---
#### log arguments
* `disable` - (Optional)(`Bool`) Disable security logging for the device.
* `event_rate` - (Optional)(`Int`) Control plane event rate (0..1500 logs per second).
* `facility_override` - (Optional)(`String`) Alternate facility for logging to remote host.
* `file` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for 'security log file' block.
  * `files` - (Optional)(`Int`) Maximum number of binary log files (2..10).
  * `name` - (Optional)(`String`) Name of binary log file.
  * `path` - (Optional)(`String`) Path to binary log files.
  * `size` - (Optional)(`Int`)  Maximum size of binary log file in megabytes (1..10).
* `format` - (Optional)(`String`) Set security log format for the device. Need to be 'binary', 'sd-syslog' or 'syslog'.
* `max_database_record` - (Optional)(`Int`) Maximum records in database (0..1000000).
* `mode` - (Optional)(`String`) Controls how security logs are processed and exported. Need to be 'event' or 'stream'.
* `rate_cap` - (Optional)(`Int`) Data plane event rate (0..5000 logs per second).
* `report` - (Optional)(`Bool`) Set security log report settings.
* `source_address` - (Optional)(`String`) Source ip address used when exporting security logs. Conflict with `source_interface`.
* `source_interface`- (Optional)(`String`) Source interface used when exporting security logs. Conflict with `source_address`.
* `transport` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for 'security log transport' block.
  * `protocol` - (Optional)(`String`) Set security log transport protocol for the device. Need to be 'tcp', 'tls' or 'udp'.
  * `tcp_connections` - (Optional)(`Int`) Set tcp connection number per-stream (1..5)
  * `tls_profile` - (Optional)(`String`) TLS profile.
* `utc_timestamp` - (Optional)(`Bool`) Use UTC time for security log timestamps.

---
#### time_wait_state arguments for tcp_session in flow
* `apply_to_half_close_state` - (Optional)(`Bool`) Apply time-wait-state timeout to half-close state.
* `session_ageout` - (Optional)(`Bool`) Allow session to ageout using service based timeout values.
* `session_timeout` - (Optional)(`Int`) Configure session timeout value for time-wait state (2..600 seconds).

---
#### feature_profile_web_filtering_juniper_enhanced_server arguments for utm
* `host` - (Optional)(`String`) Server host IP address or string host name.
* `port` - (Optional)(`Int`) Server port (1..65535).
* `proxy_profile` - (Optional)(`String`) Proxy profile.
* `routing_instance` - (Optional)(`String`) Routing instance name.

## Import

Junos security can be imported using any id, e.g.

```
$ terraform import junos_security.security random
```
