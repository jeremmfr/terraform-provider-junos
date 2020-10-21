---
layout: "junos"
page_title: "Junos: junos_system"
sidebar_current: "docs-junos-resource-system"
description: |-
  Configure static configuration in system block
---

# junos_system

-> **Note:** This resource should only create **once**. It's used to configure static (not object) options in `system` block. Destroy this resource as no effect on Junos configuration.

Configure static configuration in `system` block

## Example Usage

```hcl
# Configure system
resource junos_system "system" {
  name_server = ["192.0.2.10","192.0.2.11"]
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

* `name_server` - (Optional)(`ListOfString`) DNS name servers.
* `services` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'services' configuration.
  * `ssh` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'ssh' configuration. See the [`ssh` arguments] (#ssh-arguments) block.
* `syslog` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'syslog' configuration.
  * `archive` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once for declare 'archive' configuration. See the [`archive` arguments] (#archive-arguments) block.
  * `log_rotate_frequency` - (Optional)(`Int`) Rotate log frequency (1..59 minutes).
  * `source_address` - (Optional)(`String`) Use specified address as source address.
* `tracing_dest_override_syslog_host` - (Optional)(`String`) Send trace messages to remote syslog server.

#### ssh arguments
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

#### archive arguments
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
