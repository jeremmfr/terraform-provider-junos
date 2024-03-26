---
page_title: "Junos: junos_security"
---

# junos_security

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `security` block.  
By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `security` block

## Example Usage

```hcl
# Configure security
resource "junos_security" "security" {
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

- **clean_on_destroy** (Optional, Boolean)  
  Clean supported lines when destroy this resource.
- **alg** (Optional, Block)  
  Declare `alg` configuration.  
  See [below for nested schema](#alg-arguments).
- **flow** (Optional, Block)  
  Declare `flow` configuration.  
  See [below for nested schema](#flow-arguments).
- **forwarding_options** (Optional, Block)  
  Declare `forwarding-options` configuration.  
  See [below for nested schema](#forwarding_options-arguments).
- **forwarding_process** (Optional, Block)  
  Declare `forwarding-process` configuration.
  - **enhanced_services_mode** (Optional, Boolean)  
    Enable enhanced application services mode.
- **idp_security_package** (Optional, Block)  
  Declare `idp security-package` configuration.  
  See [below for nested schema](#idp_security_package-arguments).
- **idp_sensor_configuration** (Optional, Block)  
  Declare `idp sensor-configuration` configuration.  
  See [below for nested schema](#idp_sensor_configuration-arguments).
- **ike_traceoptions** (Optional, Block)  
  Declare `ike traceoptions` configuration.
  - **file** (Optional, Block)  
    Declare `file` configuration.  
    See [below for nested schema](#file-arguments-for-ike_traceoptions).
  - **flag** (Optional, Set of String)  
    Tracing parameters for IKE.
  - **no_remote_trace** (Optional, Boolean)  
    Disable remote tracing.
  - **rate_limit** (Optional, Number)  
    Limit the incoming rate of trace messages (0..4294967295).
- **log** (Optional, Block)  
  Declare `log` configuration.  
  See [below for nested schema](#log-arguments).
- **nat_source** (Optional, Block)  
  Declare `nat source` configuration.  
  - **address_persistent** (Optional, Boolean)  
    Allow source address to maintain same translation.
  - **interface_port_overloading_factor** (Optional, Number)  
    Port overloading factor for interface NAT.  
    Conflict with `interface_port_overloading_off`.
  - **interface_port_overloading_off** (Optional, Boolean)  
    Turn off interface port over-loading.  
    Conflict with `interface_port_overloading_factor`.
  - **pool_default_port_range** (Optional, Number)  
    Configure Source NAT default port range lower limit.  
    `pool_default_port_range_to` must also be specified.
  - **pool_default_port_range_to** (Optional, Number)  
    Configure Source NAT default port range upper limit.  
    `pool_default_port_range` must also be specified.
  - **pool_default_twin_port_range** (Optional, Number)  
    Configure Source NAT default twin port range lower limit.  
    `pool_default_twin_port_range_to` must also be specified.
  - **pool_default_twin_port_range_to** (Optional, Number)  
    Configure Source NAT default twin port range upper limit.  
    `pool_default_twin_port_range` must also be specified.
  - **pool_utilization_alarm_clear_threshold** (Optional, Number)  
    Clear threshold for pool utilization alarm (40..100).  
    `pool_utilization_alarm_raise_threshold` must also be specified.
  - **pool_utilization_alarm_raise_threshold** (Optional, Number)  
    Raise threshold for pool utilization alarm (50..100).
  - **port_randomization_disable** (Optional, Boolean)  
    Disable Source NAT port randomization.
  - **session_drop_hold_down** (Optional, Number)  
    Session drop hold down time (30..28800).
  - **session_persistence_scan** (Optional, Boolean)  
    Allow source to maintain session when session scan.
- **policies** (Optional, Block)  
  Declare `policies` configuration.
  - **policy_rematch** (Optional, Boolean)  
    Can be specified to allow session to remain open when an associated security policy is
    modified.  
    Conflict with `policy_rematch_extensive`.
  - **policy_rematch_extensive** (Optional, Boolean)  
    Can be specified to allow session to remain open when an associated security policy is modified,
    renamed, deactivated, or deleted.  
    Conflict with `policy_rematch`.
- **user_identification_auth_source** (Optional, Block)  
  Declare `user-identification authentication-source` configuration.
  - **ad_auth_priority** (Optional, Number)  
    Active-directory-authentication-table priority.  
    Larger number means lower priority, 0 for disable (0..65535).
  - **aruba_clearpass_priority** (Optional, Number)  
    ClearPass-authentication-table priority.  
    Larger number means lower priority, 0 for disable (0..65535).
  - **firewall_auth_priority** (Optional, Number)  
    Firewall-authentication priority.  
    Larger number means lower priority, 0 for disable (0..65535).
  - **local_auth_priority** (Optional, Number)  
    Local-authentication-table priority.  
    Larger number means lower priority, 0 for disable (0..65535).
  - **unified_access_control_priority** (Optional, Number)  
    Unified-access-control priority.  
    Larger number means lower priority, 0 for disable (0..65535).
- **utm** (Optional, Block)  
  Declare `utm` configuration.
  - **feature_profile_web_filtering_type** (Optional, String)  
    Configuring feature-profile web-filtering type.  
    Need to be `juniper-enhanced`, `juniper-local`, `web-filtering-none` or `websense-redirect`.
  - **feature_profile_web_filtering_juniper_enhanced_server** (Optional, Block)  
    Declare `utm feature-profile web-filtering juniper-enhanced server` configuration.  
    See [below for nested schema](#feature_profile_web_filtering_juniper_enhanced_server-arguments-for-utm).

---

### alg arguments

- **dns_disable** (Optional, Boolean)  
  Disable dns alg.
- **ftp_disable** (Optional, Boolean)  
  Disable ftp alg.
- **h323_disable** (Optional, Boolean)  
  Disable h323 alg.
- **mgcp_disable** (Optional, Boolean)  
  Disable mgcp alg.
- **msrpc_disable** (Optional, Boolean)  
  Disable msrpc alg.
- **pptp_disable** (Optional, Boolean)  
  Disable pptp alg.
- **rsh_disable** (Optional, Boolean)  
  Disable rsh alg.
- **rtsp_disable** (Optional, Boolean)  
  Disable rtsp alg.
- **sccp_disable** (Optional, Boolean)  
  Disable sccp alg.
- **sip_disable** (Optional, Boolean)  
  Disable sip alg.
- **sql_disable** (Optional, Boolean)  
  Disable sql alg.
- **sunrpc_disable** (Optional, Boolean)  
  Disable sunrpc alg.
- **talk_disable** (Optional, Boolean)  
  Disable talk alg.
- **tftp_disable** (Optional, Boolean)  
  Disable tftp alg.

---

### file arguments for ike_traceoptions

- **name** (Optional, String)  
  Name of file in which to write trace information.
- **files** (Optional, Number)  
  Maximum number of trace files (2..1000).
- **match** (Optional, String)  
  Regular expression for lines to be logged.
- **size** (Optional, Number)  
  Maximum trace file size (10240..1073741824).
- **world_readable** (Optional, Boolean)  
  Allow any user to read the log file.
- **no_world_readable** (Optional, Boolean)  
  Don't allow any user to read the log file.

---

### flow arguments

- **advanced_options** (Optional, Block)  
  Declare `flow advanced-options` configuration.
  - **drop_matching_link_local_address** (Optional, Boolean)  
    Drop matching link local address.
  - **drop_matching_reserved_ip_address** (Optional, Boolean)  
    Drop matching reserved source IP address.
  - **reverse_route_packet_mode_vr** (Optional, Boolean)  
    Allow reverse route lookup with packet mode vr.
- **aging** (Optional, Block)  
  Declare `flow aging` configuration.
  - **early_ageout** (Optional, Number)  
    Delay before device declares session invalid (1..65535 seconds).
  - **high_watermark** (Optional, Boolean)  
    Percentage of session-table capacity at which aggressive aging-out starts (0..100 percent).
  - **low_watermark** (Optional, Boolean)  
    Percentage of session-table capacity at which aggressive aging-out ends (0..100 percent).
- **allow_dns_reply** (Optional, Boolean)  
  Allow unmatched incoming DNS reply packet.
- **allow_embedded_icmp** (Optional, Boolean)  
  Allow embedded ICMP packets not matching a session to pass through.
- **allow_reverse_ecmp** (Optional, Boolean)  
  Allow reverse ECMP route lookup.
- **enable_reroute_uniform_link_check_nat** (Optional, Boolean)  
  Enable reroute check with uniform link and NAT check.
- **ethernet_switching** (Optional, Block)  
  Declare `flow ethernet-switching` configuration.
  - **block_non_ip_all** (Optional, Boolean)  
    Block all non-IP and non-ARP traffic including broadcast/multicast.
  - **bypass_non_ip_unicast** (Optional, Boolean)  
    Allow all non-IP (including unicast) traffic.
  - **bpdu_vlan_flooding** (Optional, Boolean)  
    Set 802.1D BPDU flooding based on VLAN.
  - **no_packet_flooding** (Optional, Block)  
    Stop IP flooding, send ARP/ICMP to trigger MAC learning.  
    There is one argument : **no_trace_route** (Optional, Boolean) Don't send ICMP to trigger MAC learning.
- **force_ip_reassembly** (Optional, Boolean)  
  Force to reassemble ip fragments.
- **ipsec_performance_acceleration** (Optional, Boolean)  
  Accelerate the IPSec traffic performance.
- **mcast_buffer_enhance** (Optional, Boolean)  
  Allow to hold more packets during multicast session creation.
- **pending_sess_queue_length** (Optional, String)  
  Maximum queued length per pending session.  
  Need to be `high`, `moderate` or `normal`.
- **preserve_incoming_fragment_size** (Optional, Boolean)  
  Preserve incoming fragment size for egress MTU.
- **route_change_timeout** (Optional, Number)  
  Timeout value for route change to nonexistent route (6..1800 seconds).
- **syn_flood_protection_mode** (Optional, String)  
  TCP SYN flood protection mode.  
  Need to be `syn-cookie` or `syn-proxy`.
- **sync_icmp_session** (Optional, Boolean)  
  Allow icmp sessions to sync to peer node.
- **tcp_mss** (Optional, Block)  
  Declare `flow tcp-mss` configuration.
  - **all_tcp_mss** (Optional, Number)  
    Enable MSS override for all packets with this value.
  - **gre_in** (Optional, Block)  
    Enable MSS override for all GRE packets coming out of an IPSec tunnel.  
    There is one argument : **mss** (Optional, Number) MSS Value.
  - **gre_out** (Optional, Block)  
    Enable MSS override for all GRE packets entering an IPsec tunnel.  
    There is one argument : **mss** (Optional, Number) MSS Value.
  - **ipsec_vpn** (Optional, Block)  
    Enable MSS override for all packets entering IPSec tunnel.  
    There is one argument : **mss** (Optional, Number) MSS Value.
- **tcp_session** (Optional, Block)  
  Declare `flow tcp-session` configuration.
  - **fin_invalidate_session** (Optional, Boolean)  
    Immediately end session on receipt of fin (FIN) segment.
  - **maximum_window** (Optional, String)  
    Maximum TCP proxy scaled receive window.  
    Need to be `64K`, `128K`, `256K`, `512K` or `1M`.
  - **no_sequence_check** (Optional, Boolean)  
    Disable sequence-number checking.
  - **no_syn_check** (Optional, Boolean)  
    Disable creation-time SYN-flag check.  
    Conflict with `strict_syn_check`.
  - **no_syn_check_in_tunnel** (Optional, Boolean)  
    Disable creation-time SYN-flag check for tunnel packets.  
    Conflict with `strict_syn_check`.
  - **rst_invalidate_session** (Optional, Boolean)  
    Immediately end session on receipt of reset (RST) segment.
  - **rst_sequence_check** (Optional, Boolean)  
    Check sequence number in reset (RST) segment.
  - **strict_syn_check** (Optional, Boolean)  
    Enable strict syn check.  
    Conflict with `no_sync_check` and `no_syn_check_in_tunnel`.
  - **tcp_initial_timeout** (Optional, Number)  
    Timeout for TCP session when initialization fails (4..300 seconds).
  - **time_wait_state** (Optional, Block)  
    Declare session timeout value in time-wait state.  
    See [below for nested schema](#time_wait_state-arguments-for-tcp_session-in-flow).

---

### forwarding_options arguments

- **inet6_mode** (Optional, String)  
  Forwarding mode for inet6 family.  
  Need to be `drop`, `flow-based` or `packet-based`.
- **iso_mode_packet_based** (Optional, Boolean)  
  Forwarding mode packet-based for iso family.
- **mpls_mode** (Optional, String)  
  Forwarding mode for mpls family.  
  Need to be `flow-based` or `packet-based`.

---

### idp_security_package arguments

- **automatic_enable** (Optional, Boolean)  
  Enable scheduled download and update.
- **automatic_interval** (Optional, Number)  
  Automatic interval (1..336 hours).
- **automatic_start_time** (Optional, String)  
  Automatic start time (YYYY-MM-DD.HH:MM:SS).
- **install_ignore_version_check** (Optional, Boolean)  
  Skip version check when attack database gets installed.
- **proxy_profile** (Optional, String)  
  Proxy profile of security package download.
- **source_address** (Optional, String)  
  Source address to be used for sending download request.
- **url** (Optional, String)  
  URL of Security package download.

---

### idp_sensor_configuration arguments

- **log_cache_size** (Optional, Number)  
  Log cache size (1..65535).
- **log_suppression** (Optional, Block)  
  Enable `log suppression`.
  - **disable** (Optional, Boolean)  
    Disable log suppression.
  - **include_destination_address** (Optional, Boolean)  
    Include destination address while performing a log suppression.
  - **no_include_destination_address** (Optional, Boolean)  
    Don't include destination address while performing a log suppression.
  - **max_logs_operate** (Optional, Number)  
    Maximum logs can be operate on (256..65536).
  - **max_time_report** (Optional, Number)  
    Time after suppressed logs will be reported (1..60).
  - **start_log** (Optional, Number)  
    Suppression start log (1..128).
- **packet_log** (Optional, Block)  
  Declare `packet-log` configuration.
  - **source_address** (Required, String)  
    Source IP address used to transport packetlog to a host.
  - **host_address** (Optional, String)  
    Destination host to send packetlog to.
  - **host_port** (Optional, Number)  
    Destination UDP port number (1..65536).
  - **max_sessions** (Optional, Number)  
    Max num of sessions in unit(%) (1..100).
  - **threshold_logging_interval** (Optional, Number)  
    Interval of logs for max limit session/memory reached in minutes (1..60).
  - **total_memory** (Optional, Number)  
    Total memory unit(%) (1..100).
- **security_configuration_protection_mode** (Optional, String)  
  Enable security protection mode.

---

### log arguments

- **disable** (Optional, Boolean)  
  Disable security logging for the device.
- **event_rate** (Optional, Number)  
  Control plane event rate (0..1500 logs per second).
- **facility_override** (Optional, String)  
  Alternate facility for logging to remote host.
- **file** (Optional, Block)  
  Declare `security log file` configuration.
  - **files** (Optional, Number)  
    Maximum number of binary log files (2..10).
  - **name** (Optional, String)  
    Name of binary log file.
  - **path** (Optional, String)  
    Path to binary log files.
  - **size** (Optional, Number)  
     Maximum size of binary log file in megabytes (1..10).
- **format** (Optional, String)  
  Set security log format for the device.  
  Need to be `binary`, `sd-syslog` or `syslog`.
- **max_database_record** (Optional, Number)  
  Maximum records in database (0..1000000).
- **mode** (Optional, String)  
  Controls how security logs are processed and exported.  
  Need to be `event` or `stream`.
- **rate_cap** (Optional, Number)  
  Data plane event rate (0..5000 logs per second).
- **report** (Optional, Boolean)  
  Set security log report settings.
- **source_address** (Optional, String)  
  Source ip address used when exporting security logs.  
  Conflict with `source_interface`.
- **source_interface** (Optional, String)  
  Source interface used when exporting security logs.  
  Conflict with `source_address`.
- **transport** (Optional, Block)  
  Declare `security log transport` configuration.
  - **protocol** (Optional, String)  
    Set security log transport protocol for the device.  
    Need to be `tcp`, `tls` or `udp`.
  - **tcp_connections** (Optional, Number)  
    Set tcp connection number per-stream (1..5).
  - **tls_profile** (Optional, String)  
    TLS profile.
- **utc_timestamp** (Optional, Boolean)  
  Use UTC time for security log timestamps.

---

### time_wait_state arguments for tcp_session in flow

- **apply_to_half_close_state** (Optional, Boolean)  
  Apply time-wait-state timeout to half-close state.
- **session_ageout** (Optional, Boolean)  
  Allow session to ageout using service based timeout values.
- **session_timeout** (Optional, Number)  
  Configure session timeout value for time-wait state (2..600 seconds).

---

### feature_profile_web_filtering_juniper_enhanced_server arguments for utm

- **host** (Optional, String)  
  Server host IP address or string host name.
- **port** (Optional, Number)  
  Server port (1..65535).
- **proxy_profile** (Optional, String)  
  Proxy profile.
- **routing_instance** (Optional, String)  
  Routing instance name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `security`.

## Import

Junos security can be imported using any id, e.g.

```shell
$ terraform import junos_security.security random
```
