---
layout: "junos"
page_title: "Provider: Junos"
sidebar_current: "docs-junos-index"
description: |-
  The Junos provider communicate with Junos device via netconf protocol and modify a part of configuration
---

# Junos Provider (unofficial)

The Junos provider communicate with Junos device via netconf protocol
and modify a part of configuration.

The provider allows you to manage some elements on Junos device.

You need to add netconf service:

```text
set system services netconf ssh
```

and optionally a specific user for netconf:

```text
set system login user netconf uid 200?

set system login user netconf class xxxx

set system login user netconf authentication ssh-rsa "xxxx"
```

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Junos Provider
provider "junos" {
  ip         = "${var.junos_ip_or_dns}"
  sshkeyfile = "${var.ssh_key_path}"
}

# Configure an interface
resource "junos_interface" "server1" {
  name = "ge-0/0/3"
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `ip` - (Required) This is the target for Netconf session (ip or dns name).
  It can also be sourced from the `JUNOS_HOST` environment variable.

* `sshkeyfile` - (Required) This is the path to ssh key for establish ssh
  connection. It can also be sourced from the `JUNOS_KEYFILE` environment
  variable.

* `username` - (Optional) This is the username for ssh connection.
  It can also be sourced from the `JUNOS_USERNAME` environment variable.
  Defaults to `netconf`.

* `port` - (Optional) This is the tcp port for ssh connection.
  It can also be sourced from the `JUNOS_PORT` environment variable.
  Defaults to `830`.

* `keypass` - (Optional) This is the passphrase for open ssh key file.
  It can also be sourced from the `JUNOS_KEYPASS` environment variable.
  Defaults is empty.

* `group_interface_delete` - (Optional) This is the Junos group used for remove configuration on a physical interface. See interface specifications [interface specifications](#interface-specifications). It can also be sourced from the `JUNOS_GROUP_INTERFACE_DELETE` environment variable. Default to empty.

## Interface specifications

When create a resource for a physical interface, the provider considers the interface available if there is 'apply-groups [`group_interface_delete`](#group_interface_delete)' and only this line on interface configuration.

Example if group_interface_delete => "interface-NC":

```text
ge-0/0/3 {
  apply-groups interface-NC;
}
```

When provider destroy resource for physical interface, he add this line.

If [`group_interface_delete`](#group_interface_delete) is empty the provider add this configuration on physical interface when delete resource :

```text
ge-0/0/3 {
  description NC;
  disable;
}
```

and considers the interface available if the is this lines and only this lines on interface.

## Environment variables

You can export the `TFJUNOS_LOG_PATH` environment variable for a more detailed log (netconf) in the specified file.

You can export the `TFJUNOS_SLEEP` environment variable to change the number of seconds of standby while waiting for Terraform to lock candidate configuration on a Junos device. The default value is `10`.

You can export the `TFJUNOS_SLEEP_SHORT` environment variable to change the number of milliseconds to wait after Terraform executes an action on the Junos device. Default to `100`
