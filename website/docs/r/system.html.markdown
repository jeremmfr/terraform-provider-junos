---
layout: "junos"
page_title: "Junos: junos_system"
sidebar_current: "docs-junos-resource-system"
description: |-
  Configure static configuration in system block
---

# junos_system

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `system` block. Destroy this resource has no effect on the Junos configuration.  
There is an exception for `system root-authentication` static block. It's can be configured with the dedicated `junos_system_root_authentication` resource.

Configure static configuration in `system` block (except `system root-authentication` block)

## Example Usage

```hcl
# Configure system
resource junos_system "system" {
  host_name   = "MyJunOS-device"
  name_server = ["192.0.2.10", "192.0.2.11"]
  services {
    ssh {
      root_login = "deny"
    }
  }
  syslog {
    archive {
      files          = 5
      world_readable = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `authentication_order` - (Optional)(`ListOfString`) Order in which authentication methods are invoked.
* `auto_snapshot` - (Optional)(`Bool`) Enable auto-snapshot when boots from alternate slice.
* `default_address_selection` - (Optional)(`Bool`) Use loopback interface as source address for locally generated packets.
* `domain_name` - (Optional)(`String`) Domain name.
* `host_name` - (Optional)(`String`) Hostname.
* `inet6_backup_router` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'inet6-backup-router' configuration.
  * `address` - (Optional)(`String`) Address of router to use while booting.
  * `destination` - (Optional)(`ListOfString`) Destination networks reachable through the router.
* `internet_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'internet-options' configuration. See the [`internet_options` arguments] (#internet_options-arguments) block.
* `login` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'login' configuration. See the [`login` arguments] (#login-arguments) block.
* `max_configuration_rollbacks` - (Optional)(`Int`) Maximum rollback configuration (0..49).
* `max_configurations_on_flash` - (Optional)(`Int`) Number of configuration files stored on flash (0..49).
* `name_server` - (Optional)(`ListOfString`) DNS name servers.
* `no_multicast_echo` - (Optional)(`Bool`) Disable responding to ICMP echo requests sent to multicast group addresses.
* `no_ping_record_route` - (Optional)(`Bool`) Do not insert IP address in ping replies.
* `no_ping_time_stamp` - (Optional)(`Bool`) Do not insert time stamp in ping replies.
* `no_redirects` - (Optional)(`Bool`) Disable ICMP redirects.
* `no_redirects_ipv6` - (Optional)(`Bool`) Disable IPV6 ICMP redirects.
* `services` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'services' configuration.
  * `ssh` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ssh' configuration. See the [`ssh` arguments for services] (#ssh-arguments-for-services) block.
* `syslog` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'syslog' configuration.
  * `archive` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'archive' configuration. See the [`archive` arguments for syslog] (#archive-arguments-for-syslog) block.
  * `log_rotate_frequency` - (Optional)(`Int`) Rotate log frequency (1..59 minutes).
  * `source_address` - (Optional)(`String`) Use specified address as source address.
* `time_zone` - (Optional)(`String`) Time zone name or POSIX-compliant time zone string (<continent>/<major-city> or <time-zone>).
* `tracing_dest_override_syslog_host` - (Optional)(`String`) Send trace messages to remote syslog server.

---
#### internet_options arguments
* `gre_path_mtu_discovery` - (Optional)(`Bool`) Enable path MTU discovery for GRE tunnels. Conflict with `no_gre_path_mtu_discovery`.
* `icmpv4_rate_limit` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'icmpv4-rate-limit' configuration.
  * `bucket_size` - (Optional)(`Int`) ICMP rate-limiting maximum bucket size (seconds).
  * `packet-rate` - (Optional)(`Int`) ICMP rate-limiting packets earned per second.
* `icmpv6_rate_limit` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'icmpv6-rate-limit' configuration.
  * `bucket_size` - (Optional)(`Int`) ICMPv6 rate-limiting maximum bucket size (seconds).
  * `packet-rate` - (Optional)(`Int`) ICMPv6 rate-limiting packets earned per second.
* `ipip_path_mtu_discovery` - (Optional)(`Bool`) Enable path MTU discovery for IP-IP tunnels. Conflict with `no_ipip_path_mtu_discovery`.
* `ipv6_duplicate_addr_detection_transmits` - (Optional)(`Int`) IPv6 Duplicate address detection transmits (0..20).
* `ipv6_path_mtu_discovery` - (Optional)(`Bool`) Enable IPv6 Path MTU discovery. Conflict with `no_ipv6_path_mtu_discovery`.
* `ipv6_path_mtu_discovery_timeout` - (Optional)(`Int`) IPv6 Path MTU Discovery timeout (5..71582788 minutes).
* `ipv6_reject_zero_hop_limit` - (Optional)(`Bool`) Enable dropping IPv6 packets with zero hop-limit. Conflict with `no_ipv6_reject_zero_hop_limit`.
* `no_gre_path_mtu_discovery` - (Optional)(`Bool`) Don't enable path MTU discovery for GRE tunnels. Conflict with `gre_path_mtu_discovery`.
* `no_ipip_path_mtu_discovery` - (Optional)(`Bool`) Don't enable path MTU discovery for IP-IP tunnels. Conflict with `ipip_path_mtu_discovery`.
* `no_ipv6_path_mtu_discovery` - (Optional)(`Bool`) Don't enable IPv6 Path MTU discovery. Conflict with `ipv6_path_mtu_discovery`.
* `no_ipv6_reject_zero_hop_limit` - (Optional)(`Bool`) Don't enable dropping IPv6 packets with zero hop-limit. Conflict with `ipv6_reject_zero_hop_limit`.
* `no_path_mtu_discovery` - (Optional)(`Bool`) Don't enable Path MTU discovery on TCP connections. Conflict with `path_mtu_discovery`.
* `no_source_quench` - (Optional)(`Bool`) Don't react to incoming ICMP Source Quench messages. Conflict with `source_quench`
* `no_tcp_reset` - (Optional)(`String`) Do not send RST TCP packet for packets sent to non-listening ports. Need to be `drop-all-tcp` or `drop-tcp-with-syn-only`.
* `no_tcp_rfc1323` - (Optional)(`Bool`) Disable RFC 1323 TCP extensions.
* `no_tcp_rfc1323_paws` - (Optional)(`Bool`) Disable RFC 1323 Protection Against Wrapped Sequence Number extension.
* `path_mtu_discovery` - (Optional)(`Bool`) Enable Path MTU discovery on TCP connections. Conflict with `no_path_mtu_discovery`.
* `source_port_upper_limit` - (Optional)(`Int`) Specify upper limit of source port selection range (5000..65535).
* `source_quench` - (Optional)(`Bool`) React to incoming ICMP Source Quench messages. Conflict with `no_source_quench`.
* `tcp_drop_synfin_set` - (Optional)(`Bool`) Drop TCP packets that have both SYN and FIN flags.
* `tcp_mss` - (Optional)(`Int`) Maximum value of TCP MSS for IPV4 traffic (64..65535 bytes).

---
#### login arguments
* `announcement` - (Optional)(`String`) System announcement message (displayed after login).
* `deny_sources_address` - (Optional)(`ListOfString`) Sources from which logins are denied.
* `idle_timeout` - (Optional)(`Int`) Maximum idle time before logout (1..60 minutes).
* `message` - (Optional)(`String`) System login message.
* `password` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'password' configuration.
  * `change_type` - (Optional)(`String`) Password change type.
  * `format` - (Optional)(`String`) Encryption method to use for password.
  * `maximum_length` - (Optional)(`Int`) Maximum password length for all users (20..128).
  * `minimum_changes` - (Optional)(`Int`) Minimum number of changes in password (1..128).
  * `minimum_character_changes` - (Optional)(`Int`) Minimum number of character changes between old and new passwords (4..15).
  * `minimum_length` - (Optional)(`Int`) Minimum password length for all users (6..20).
  * `minimum_lower_cases` - (Optional)(`Int`) Minimum number of lower-case class characters in password (1..128)
  * `minimum_numerics` - (Optional)(`Int`) Minimum number of numeric class characters in password (1..128).
  * `minimum_punctuations` - (Optional)(`Int`) Minimum number of punctuation class characters in password (1..128).
  * `minimum_reuse` - (Optional)(`Int`) Minimum number of old passwords which should not be same as the new password (1..20).
  * `minimum_upper_cases` - (Optional)(`Int`) Minimum number of upper-case class characters in password (1..128).
* `retry_options` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'retry-options' configuration.
  * `backoff_factor` - (Optional)(`Int`) Delay factor after 'backoff-threshold' password failures (5..10).
  * `backoff_threshold` - (Optional)(`Int`) Number of password failures before delay is introduced (1..3).
  * `lockout_period` - (Optional)(`Int`) Amount of time user account is locked after 'tries-before-disconnect' failures (1..43200 minutes).
  * `maximum_time` - (Optional)(`Int`) Maximum time the connection will remain for user to enter username and password (20..300).
  * `minimum_time` - (Optional)(`Int`) Minimum total connection time if all attempts fail (20..60).
  * `tries_before_disconnect` - (Optional)(`Int`) Number of times user is allowed to try password (2..10).

---
#### ssh arguments for services
* `authentication_order` - (Optional)(`ListOfString`) Order in which authentication methods are invoked.
* `ciphers` - (Optional)(`ListOfString`) Specify the ciphers allowed for protocol version 2.
* `client_alive_count_max` - (Optional)(`Int`) Threshold of missing client-alive responses that triggers a disconnect (0..255).
* `client_alive_interval` - (Optional)(`Int`) Frequency of client-alive requests (0..65535 seconds).
* `connection_limit` - (Optional)(`Int`) Maximum number of allowed connections (1..250).
* `fingerprint_hash` - (Optional)(`String`) Configure hash algorithm used when displaying key fingerprints.
* `hostkey_algorithm` - (Optional)(`ListOfString`) Specify permissible SSH host-key algorithms.
* `key_exchange` - (Optional)(`ListOfString`) Specify ssh key-exchange for Diffie-Hellman keys.
* `macs` - (Optional)(`ListOfString`) Message Authentication Code algorithms allowed (SSHv2).
* `max_pre_authentication_packets` - (Optional)(`Int`) Maximum number of pre-authentication SSH packets per single SSH connection (20..2147483647).
* `max_sessions_per_connection` - (Optional)(`Int`) Maximum number of sessions per single SSH connection (1..65535).
* `no_passwords` - (Optional)(`Bool`) Disables ssh password based authentication.
* `no_public_keys` - (Optional)(`Bool`) Disables ssh public key based authentication.
* `port` - (Optional)(`Int`) Port number to accept incoming connections (1..65535).
* `protocol_version` - (Optional)(`ListOfString`) Specify ssh protocol versions supported.
* `rate_limit` - (Optional)(`Int`) Maximum number of connections per minute (1..250).
* `root_login` - (Optional)(`String`) Configure root access via ssh. Need to be 'allow', 'deny' or 'deny-password'.
* `no_tcp_forwarding` - (Optional)(`Bool`) Do not allow forwarding TCP connections via SSH.
* `tcp_forwarding` - (Optional)(`Bool`) Allow forwarding TCP connections via SSH.

---
#### archive arguments for syslog
* `binary_data` - (Optional)(`Bool`) Mark file as if it contains binary data.
* `no_binary_data` - (Optional)(`Bool`) Don't mark file as if it contains binary data.
* `files` - (Optional)(`Int`) Number of files to be archived (1..1000).
* `size` - (Optional)(`Int`) Size of files to be archived (65536..1073741824 bytes)
* `world_readable` - (Optional)(`Bool`) Allow any user to read the log file.
* `no_world_readable` - (Optional)(`Bool`) Don't allow any user to read the log file.

## Import

Junos system can be imported using any id, e.g.

```
$ terraform import junos_system.system random
```
