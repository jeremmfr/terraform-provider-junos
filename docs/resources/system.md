---
page_title: "Junos: junos_system"
---

# junos_system

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `system` block.  
There is an exception for `system root-authentication` static block.
It's can be configured with the dedicated `junos_system_root_authentication` resource.  
Destroy this resource has no effect on the Junos configuration.

Configure static configuration in `system` block (except `system root-authentication` block)

## Example Usage

```hcl
# Configure system
resource "junos_system" "system" {
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

- **accounting** (Optional, Block)  
  Declare `accounting` configuration.  
  - **events** (Required, Set of String)  
    Events to be logged.
  - **destination_radius** (Optional, Boolean)  
    Send RADIUS accounting records.
  - **destination_radius_server** (Optional, Block List)  
    For each address, RADIUS accounting server configuration.  
    See [below for nested schema](#destination_radius_server-arguments-for-accounting).
  - **destination_tacplus** (Optional, Boolean)  
    Send TACACS+ accounting records.
  - **destination_tacplus_server** (Optional, Block List)  
    For each address, TACACS+ accounting server configuration.  
    See [below for nested schema](#destination_tacplus_server-arguments-for-accounting).
  - **enhanced_avs_max** (Optional, Number)  
    No. of AV pairs each of which can store a max of 250 Bytes.
- **archival_configuration** (Optional, Block)  
  Declare `archival configuration` configuration.  
  See [below for nested schema](#archival_configuration-arguments).
- **authentication_order** (Optional, List of String)  
  Order in which authentication methods are invoked.  
  Element need to be `password`, `radius` or `tacplus`.
- **auto_snapshot** (Optional, Boolean)  
  Enable auto-snapshot when boots from alternate slice.
- **default_address_selection** (Optional, Boolean)  
  Use loopback interface as source address for locally generated packets.
- **domain_name** (Optional, String)  
  Domain name.
- **host_name** (Optional, String)  
  Hostname.
- **inet6_backup_router** (Optional, Block)  
  Declare `inet6-backup-router` configuration.
  - **address** (Optional, String)  
    Address of router to use while booting.
  - **destination** (Optional, Set of String)  
    Destination networks reachable through the router.
- **internet_options** (Optional, Block)  
  Declare `internet-options` configuration.  
  See [below for nested schema](#internet_options-arguments).
- **license** (Optional, Block)  
  Declare `license` configuration.
  - **autoupdate** (Optional, Boolean)  
    Enable autoupdate license keys.
  - **autoupdate_password** (Optional, String, Sensitive)  
    Password for autoupdate license keys from license servers.  
    `autoupdate_url` needs to be set.  
  - **autoupdate_url** (Optional, String)  
    Url for autoupdate license keys from license servers.  
    `autoupdate` needs to be set.
  - **renew_before_expiration** (Optional, Number)  
    License renewal lead time before expiration, in days (0..60).  
    `renew_interval` needs to be set.
  - **renew_interval** (Optional, Number)  
    License checking interval, in hours (1..336).  
    `renew_before_expiration` needs to be set.
- **login** (Optional, Block)  
  Declare `login` configuration.  
  See [below for nested schema](#login-arguments).
- **max_configuration_rollbacks** (Optional, Number)  
  Maximum rollback configuration (0..49).
- **max_configurations_on_flash** (Optional, Number)  
  Number of configuration files stored on flash (0..49).
- **name_server** (Optional, List of String)  
  DNS name servers.  
  Conflict with `name_server_opts`.
- **name_server_opts** (Optional, Block List)  
  DNS name servers with optional options.  
  Conflict with `name_server`.
  - **address** (Required, String)  
    Address of the name server.
  - **routing_instance** (Optional, String)  
    Routing instance through which the name server is reachable.
- **no_multicast_echo** (Optional, Boolean)  
  Disable responding to ICMP echo requests sent to multicast group addresses.
- **no_ping_record_route** (Optional, Boolean)  
  Do not insert IP address in ping replies.
- **no_ping_time_stamp** (Optional, Boolean)  
  Do not insert time stamp in ping replies.
- **no_redirects** (Optional, Boolean)  
  Disable ICMP redirects.
- **no_redirects_ipv6** (Optional, Boolean)  
  Disable IPV6 ICMP redirects.
- **ntp** (Optional, Block)  
  Declare `ntp` configuration.
  - **boot_server** (Optional, String)  
    Server to query during boot sequence.
  - **broadcast_client** (Optional, Boolean)  
    Listen to broadcast NTP.
  - **interval_range** (Optional, Number)  
    Set the minpoll and maxpoll interval range (0..3).
  - **multicast_client** (Optional, Boolean)  
    Listen to multicast NTP.
  - **multicast_client_address** (Optional, String)  
    Multicast address to listen to.  
    `multicast_client` need to be set to true.
  - **threshold_action** (Optional, String)  
    Select actions for NTP abnormal adjustment.  
    Need to be `accept` or `reject`.  
    `threshold_value` needs to be set.
  - **threshold_value** (Optional, Number)  
    Set the maximum threshold(sec) allowed for NTP adjustment (1..600).
    `threshold_action` needs to be set.
- **ports** (Optional, Block)  
  Declare `ports` configuration.
  - **auxiliary_authentication_order** (Optional, List of String)  
    Order in which authentication methods are invoked on auxiliary port.  
    Element need to be `password`, `radius` or `tacplus`.
  - **auxiliary_disable** (Optional, Boolean)  
    Disable console on auxiliary port.
  - **auxiliary_insecure** (Optional, Boolean)  
    Disallow superuser access on auxiliary port.
  - **auxiliary_logout_on_disconnect** (Optional, Boolean)  
    Log out the console session when cable is unplugged.
  - **auxiliary_type** (Optional, String)  
    Terminal type on auxiliary port.
  - **console_authentication_order** (Optional, List of String)  
    Order in which authentication methods are invoked on console port.  
    Element need to be `password`, `radius` or `tacplus`.
  - **console_disable** (Optional, Boolean)  
    Disable console on console port.
  - **console_insecure** (Optional, Boolean)  
    Disallow superuser access on console port.
  - **console_logout_on_disconnect** (Optional, Boolean)  
    Log out the console session when cable is unplugged.
  - **console_type** (Optional, String)  
    Terminal type on console port.
- **radius_options_attributes_nas_id** (Optional, String)  
  Value of NAS-ID in outgoing RADIUS packets.
- **radius_options_attributes_nas_ipaddress** (Optional, String)  
  Value of NAS-IP-Address in outgoing RADIUS packets.
- **radius_options_enhanced_accounting** (Optional, Boolean)  
  Include authentication method, remote port and user-privileges in `login` accounting.
- **radius_options_password_protocol_mschapv2** (Optional, Boolean)  
  MSCHAP version 2 for password protocol used in RADIUS packets.
- **services** (Optional, Block)  
  Declare `services` configuration.
  - **netconf_ssh** (Optional, Block)  
    Declare `netconf ssh` configuration.  
    See [below for nested schema](#netconf_ssh-arguments-for-services).
  - **netconf_traceoptions** (Optional, Block)  
    Declare `netconf traceoptions` configuration.  
    See [below for nested schema](#netconf_traceoptions-arguments-for-services).
  - **ssh** (Optional, Block)  
    Declare `ssh` configuration.  
    See [below for nested schema](#ssh-arguments-for-services).
  - **web_management_http** (Optional, Block)  
    Enable `web-management http`.  
    See [below for nested schema](#web_management_http-arguments-for-services).
  - **web_management_https** (Optional, Block)  
    Declare `web-management https` configuration.  
    See [below for nested schema](#web_management_https-arguments-for-services).
  - **web_management_session_idle_timeout** (Optional, Number)  
    Default timeout of web-management sessions (1..1440 minutes).
  - **web_management_session_limit** (Optional, Number)  
    Maximum number of web-management sessions to allow (1..1024).
- **syslog** (Optional, Block)  
  Declare `syslog` configuration.
  - **archive** (Optional, Block)  
    Declare `archive` configuration.  
    See [below for nested schema](#archive-arguments-for-syslog).
  - **console** (Optional, Block)  
    Declare `console` configuration.  
    See [below for nested schema](#console-arguments-for-syslog).
  - **log_rotate_frequency** (Optional, Number)  
    Rotate log frequency (1..59 minutes).
  - **source_address** (Optional, String)  
    Use specified address as source address.
  - **time_format_millisecond** (Optional, Boolean)  
    Include milliseconds in system log timestamp.
  - **time_format_year** (Optional, Boolean)  
    Include year in system log timestamp.
- **tacplus_options_authorization_time_interval** (Optional, Number)  
  TACACS+ authorization refresh time interval (15..1440 minutes).
- **tacplus_options_enhanced_accounting** (Optional, Boolean)  
  Include authentication method, remote port and user-privileges in `login` accounting.
- **tacplus_options_exclude_cmd_attribute** (Optional, Boolean)  
  In start/stop requests, do not include `cmd` attribute.  
  Conflict with `tacplus_options_no_cmd_attribute_value`.
- **tacplus_options_no_cmd_attribute_value** (Optional, Boolean)  
  In start/stop requests, set `cmd` attribute value to empty string.  
  Conflict with `tacplus_options_exclude_cmd_attribute`.
- **tacplus_options_service_name** (Optional, String)  
  TACACS+ service name.
- **tacplus_options_strict_authorization** (Optional, Boolean)  
  Deny login if authorization request fails.  
  Conflict with `tacplus_options_no_strict_authorization`.
- **tacplus_options_no_strict_authorization** (Optional, Boolean)  
  Don't deny login if authorization request fails.  
  Conflict with `tacplus_options_strict_authorization`.
- **tacplus_options_timestamp_and_timezone** (Optional, Boolean)  
  In start/stop accounting packets, include `start-time`, `stop-time` and `timezone` attributes.
- **time_zone** (Optional, String)  
  Time zone name or POSIX-compliant time zone string (`<continent>`/`<major-city>` or `<time-zone>`).
- **tracing_dest_override_syslog_host** (Optional, String)  
  Send trace messages to remote syslog server.

---

### destination_radius_server arguments for accounting

- **address** (Required, String)  
  RADIUS server address.
- **secret** (Required, String, Sensitive)  
  Shared secret with the RADIUS server.
- **accounting_port** (Optional, Number)  
  RADIUS server accounting port number (1..65535).
- **accounting_retry** (Optional, Number)  
  Accounting retry attempts (0..100).
- **accounting_timeout** (Optional, Number)  
  Accounting request timeout period (0..1000 seconds).
- **dynamic_request_port** (Optional, Number)  
  RADIUS client dynamic request port number (1..65535).
- **max_outstanding_requests** (Optional, Number)  
  Maximum requests in flight to server (0..2000).
- **port** (Optional, Number)  
  RADIUS server authentication port number (1..65535).
- **preauthentication_port** (Optional, Number)  
  RADIUS server preauthentication port number (1..65535).
- **preauthentication_secret** (Optional, String, Sensitive)  
  Preauthentication shared secret with the RADIUS server.
- **retry** (Optional, Number)  
  Retry attempts (1..100).
- **routing_instance** (Optional, String)  
  Routing instance.
- **source_address** (Optional, String)  
  Use specified address as source address.
- **timeout** (Optional, Number)  
  Request timeout period (1..1000 seconds).

---

### destination_tacplus_server arguments for accounting

- **address** (Required, String)  
  TACACS+ authentication server address.
- **port** (Optional, Number)  
  TACACS+ authentication server port number (1..65535).
- **routing_instance** (Optional, String)  
  Routing instance.
- **secret** (Optional, String, Sensitive)  
  Shared secret with the authentication server.
- **single_connection** (Optional, Boolean)  
  Optimize TCP connection attempts.
- **source_address** (Optional, String)  
  Use specified address as source address.
- **timeout** (Optional, Number)  
  Request timeout period (1..90 seconds).

---

### archival_configuration arguments

- **archive_site** (Required, Block List)  
  For each url, configure archive-site destination.
  - **url** (Required, String)  
    URLs to receive configuration files.
  - **password** (Optional, String, Sensitive)  
    Password for login into the archive site.  
- **transfer_interval** (Optional, Number)  
  Frequency at which file transfer happens (15..2880 minutes).  
  Need to set one of `transfer_interval` or `transfer_on_commit`.
- **transfer_on_commit** (Optional, Boolean)  
  Transfer after each commit.  
  Need to set one of `transfer_interval` or `transfer_on_commit`.

---

### internet_options arguments

- **gre_path_mtu_discovery** (Optional, Boolean)  
  Enable path MTU discovery for GRE tunnels.  
  Conflict with `no_gre_path_mtu_discovery`.
- **no_gre_path_mtu_discovery** (Optional, Boolean)  
  Don't enable path MTU discovery for GRE tunnels.  
  Conflict with `gre_path_mtu_discovery`.
- **icmpv4_rate_limit** (Optional, Block)  
  Declare `icmpv4-rate-limit` configuration.
  - **bucket_size** (Optional, Number)  
    ICMP rate-limiting maximum bucket size (seconds).
  - **packet-rate** (Optional, Number)  
    ICMP rate-limiting packets earned per second.
- **icmpv6_rate_limit** (Optional, Block)  
  Declare `icmpv6-rate-limit` configuration.
  - **bucket_size** (Optional, Number)  
    ICMPv6 rate-limiting maximum bucket size (seconds).
  - **packet-rate** (Optional, Number)  
    ICMPv6 rate-limiting packets earned per second.
- **ipip_path_mtu_discovery** (Optional, Boolean)  
  Enable path MTU discovery for IP-IP tunnels.  
  Conflict with `no_ipip_path_mtu_discovery`.
- **no_ipip_path_mtu_discovery** (Optional, Boolean)  
  Don't enable path MTU discovery for IP-IP tunnels.  
  Conflict with `ipip_path_mtu_discovery`.
- **ipv6_duplicate_addr_detection_transmits** (Optional, Number)  
  IPv6 Duplicate address detection transmits (0..20).
- **ipv6_path_mtu_discovery** (Optional, Boolean)  
  Enable IPv6 Path MTU discovery.  
  Conflict with `no_ipv6_path_mtu_discovery`.
- **no_ipv6_path_mtu_discovery** (Optional, Boolean)  
  Don't enable IPv6 Path MTU discovery.  
  Conflict with `ipv6_path_mtu_discovery`.
- **ipv6_path_mtu_discovery_timeout** (Optional, Number)  
  IPv6 Path MTU Discovery timeout (5..71582788 minutes).
- **ipv6_reject_zero_hop_limit** (Optional, Boolean)  
  Enable dropping IPv6 packets with zero hop-limit.  
  Conflict with `no_ipv6_reject_zero_hop_limit`.
- **no_ipv6_reject_zero_hop_limit** (Optional, Boolean)  
  Don't enable dropping IPv6 packets with zero hop-limit.  
  Conflict with `ipv6_reject_zero_hop_limit`.
- **no_tcp_reset** (Optional, String)  
  Do not send RST TCP packet for packets sent to non-listening ports.  
  Need to be `drop-all-tcp` or `drop-tcp-with-syn-only`.
- **no_tcp_rfc1323** (Optional, Boolean)  
  Disable RFC 1323 TCP extensions.
- **no_tcp_rfc1323_paws** (Optional, Boolean)  
  Disable RFC 1323 Protection Against Wrapped Sequence Number extension.
- **path_mtu_discovery** (Optional, Boolean)  
  Enable Path MTU discovery on TCP connections.  
  Conflict with `no_path_mtu_discovery`.
- **no_path_mtu_discovery** (Optional, Boolean)  
  Don't enable Path MTU discovery on TCP connections.  
  Conflict with `path_mtu_discovery`.
- **source_port_upper_limit** (Optional, Number)  
  Specify upper limit of source port selection range (5000..65535).
- **source_quench** (Optional, Boolean)  
  React to incoming ICMP Source Quench messages.  
  Conflict with `no_source_quench`.
- **no_source_quench** (Optional, Boolean)  
  Don't react to incoming ICMP Source Quench messages.  
  Conflict with `source_quench`
- **tcp_drop_synfin_set** (Optional, Boolean)  
  Drop TCP packets that have both SYN and FIN flags.
- **tcp_mss** (Optional, Number)  
  Maximum value of TCP MSS for IPV4 traffic (64..65535 bytes).

---

### login arguments

- **announcement** (Optional, String)  
  System announcement message (displayed after login).
- **deny_sources_address** (Optional, Set of String)  
  Sources from which logins are denied.
- **idle_timeout** (Optional, Number)  
  Maximum idle time before logout (1..60 minutes).
- **message** (Optional, String)  
  System login message.
- **password** (Optional, Block)  
  Declare `password` configuration.
  - **change_type** (Optional, String)  
    Password change type.
  - **format** (Optional, String)  
    Encryption method to use for password.
  - **maximum_length** (Optional, Number)  
    Maximum password length for all users (20..128).
  - **minimum_changes** (Optional, Number)  
    Minimum number of changes in password (1..128).
  - **minimum_character_changes** (Optional, Number)  
    Minimum number of character changes between old and new passwords (4..15).
  - **minimum_length** (Optional, Number)  
    Minimum password length for all users (6..20).
  - **minimum_lower_cases** (Optional, Number)  
    Minimum number of lower-case class characters in password (1..128)
  - **minimum_numerics** (Optional, Number)  
    Minimum number of numeric class characters in password (1..128).
  - **minimum_punctuations** (Optional, Number)  
    Minimum number of punctuation class characters in password (1..128).
  - **minimum_reuse** (Optional, Number)  
    Minimum number of old passwords which should not be same as the new password (1..20).
  - **minimum_upper_cases** (Optional, Number)  
    Minimum number of upper-case class characters in password (1..128).
- **retry_options** (Optional, Block)  
  Declare `retry-options` configuration.
  - **backoff_factor** (Optional, Number)  
    Delay factor after `backoff-threshold` password failures (5..10).
  - **backoff_threshold** (Optional, Number)  
    Number of password failures before delay is introduced (1..3).
  - **lockout_period** (Optional, Number)  
    Amount of time user account is locked after `tries-before-disconnect` failures (1..43200 minutes).
  - **maximum_time** (Optional, Number)  
    Maximum time the connection will remain for user to enter username and password (20..300).
  - **minimum_time** (Optional, Number)  
    Minimum total connection time if all attempts fail (20..60).
  - **tries_before_disconnect** (Optional, Number)  
    Number of times user is allowed to try password (2..10).

---

### netconf_ssh arguments for services

- **client_alive_count_max** (Optional, Number)  
  Threshold of missing client-alive responses that triggers a disconnect (0..255).
- **client_alive_interval** (Optional, Number)  
  Frequency of client-alive requests (0..65535 seconds).
- **connection_limit** (Optional, Number)  
  Limit number of simultaneous connections (1..250 connections).
- **rate_limit** (Optional, Number)  
  Limit incoming connection rate (1..250 connections per minute).

---

### netconf_traceoptions arguments for services

- **file_name** (Optional, String)  
  Name of file in which to write trace information.
- **file_files** (Optional, Number)  
  Maximum number of trace files (2..1000).
- **file_match** (Optional, String)  
  Regular expression for lines to be logged.
- **file_size** (Optional, Number)  
  Maximum trace file size (10240..1073741824).
- **file_world_readable** (Optional, Boolean)  
  Allow any user to read the log file.
- **file_no_world_readable** (Optional, Boolean)  
  Don't allow any user to read the log file.
- **flag** (Optional, Set of String)  
  Tracing parameters.  
  Element need to be `all`, `debug`, `incoming` or `outgoing`.
- **no_remote_trace** (Optional, Boolean)  
  Disable remote tracing.
- **on_demand** (Optional, Boolean)  
  Enable on-demand tracing.

---

### ssh arguments for services

- **authentication_order** (Optional, List of String)  
  Order in which authentication methods are invoked.  
  Element need to be `password`, `radius` or `tacplus`.
- **ciphers** (Optional, Set of String)  
  Specify the ciphers allowed for protocol version 2.
- **client_alive_count_max** (Optional, Number)  
  Threshold of missing client-alive responses that triggers a disconnect (0..255).
- **client_alive_interval** (Optional, Number)  
  Frequency of client-alive requests (0..65535 seconds).
- **connection_limit** (Optional, Number)  
  Maximum number of allowed connections (1..250).
- **fingerprint_hash** (Optional, String)  
  Configure hash algorithm used when displaying key fingerprints.
- **hostkey_algorithm** (Optional, Set of String)  
  Specify permissible SSH host-key algorithms.
- **key_exchange** (Optional, Set of String)  
  Specify ssh key-exchange for Diffie-Hellman keys.
- **log_key_changes** (Optional, Boolean)  
  Log changes to authorized keys to syslog.
- **macs** (Optional, Set of String)  
  Message Authentication Code algorithms allowed (SSHv2).
- **max_pre_authentication_packets** (Optional, Number)  
  Maximum number of pre-authentication SSH packets per single SSH connection (20..2147483647).
- **max_sessions_per_connection** (Optional, Number)  
  Maximum number of sessions per single SSH connection (1..65535).
- **no_passwords** (Optional, Boolean)  
  Disables ssh password based authentication.
- **no_public_keys** (Optional, Boolean)  
  Disables ssh public key based authentication.
- **port** (Optional, Number)  
  Port number to accept incoming connections (1..65535).
- **protocol_version** (Optional, Set of String)  
  Specify ssh protocol versions supported.
  Element need to be `v1` or `v2`.
- **rate_limit** (Optional, Number)  
  Maximum number of connections per minute (1..250).
- **root_login** (Optional, String)  
  Configure root access via ssh.  
  Need to be `allow`, `deny` or `deny-password`.
- **tcp_forwarding** (Optional, Boolean)  
  Allow forwarding TCP connections via SSH.
- **no_tcp_forwarding** (Optional, Boolean)  
  Do not allow forwarding TCP connections via SSH.

---

### web_management_http arguments for services

- **interface** (Optional, Set of String)  
  Specify the name of one or more interfaces.
- **port** (Optional, Number)  
  Port number to connect to HTTP service (1..65535).

---

### web_management_https arguments for services

-> **Note:** One of `local_certificate`, `pki_local_certificate` or `system_generated_certificate`
arguments is required.

- **interface** (Optional, Set of String)  
  Specify the name of one or more interfaces.
- **local_certificate** (Optional, String)  
  Specify the name of the certificate.
- **pki_local_certificate** (Optional, String)  
  Specify the name of the certificate that is generated by the PKI and authenticated by a CA.
- **port** (Optional, Number)  
  Port number to connect to HTTPS service (1..65535).
- **system_generated_certificate** (Optional, Boolean)  
  Will automatically generate a self-signed certificate.

---

### archive arguments for syslog

- **binary_data** (Optional, Boolean)  
  Mark file as if it contains binary data.
- **no_binary_data** (Optional, Boolean)  
  Don't mark file as if it contains binary data.
- **files** (Optional, Number)  
  Number of files to be archived (1..1000).
- **size** (Optional, Number)  
  Size of files to be archived (65536..1073741824 bytes)
- **world_readable** (Optional, Boolean)  
  Allow any user to read the log file.
- **no_world_readable** (Optional, Boolean)  
  Don't allow any user to read the log file.

---

### console arguments for syslog

- **any_severity** (Optional, String)  
  All facilities severity.
- **authorization_severity** (Optional, String)  
  Authorization system severity.
- **changelog_severity** (Optional, String)  
  Configuration change log severity.
- **conflictlog_severity** (Optional, String)  
  Configuration conflict log severity.
- **daemon_severity** (Optional, String)  
  Various system processes severity.
- **dfc_severity** (Optional, String)  
  Dynamic flow capture severity.
- **external_severity** (Optional, String)  
  Local external applications severity.
- **firewall_severity** (Optional, String)  
  Firewall filtering system severity.
- **ftp_severity** (Optional, String)  
  FTP process severity.
- **interactivecommands_severity** (Optional, String)  
  Commands executed by the UI severity.
- **kernel_severity** (Optional, String)  
  Kernel severity.
- **ntp_severity** (Optional, String)  
  NTP process severity.
- **pfe_severity** (Optional, String)  
  Packet Forwarding Engine severity.
- **security_severity** (Optional, String)  
  Security related severity.
- **user_severity** (Optional, String)  
  User processes severity.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `system`.

## Import

Junos system can be imported using any id, e.g.

```shell
$ terraform import junos_system.system random
```
