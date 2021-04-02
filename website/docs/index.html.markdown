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

## Provider installation
For manual installation see [README on github](https://github.com/jeremmfr/terraform-provider-junos#install-binary-on-disk)<br>
For automatic installation (terraform >= 0.13.x) use [registry](https://registry.terraform.io/providers/jeremmfr/junos/) :

```hcl
terraform {
  required_providers {
    junos = {
      source = "jeremmfr/junos"
    }
  }
}
```

## Configure netconf

You need to add netconf service on your Junos device:

```text
set system services netconf ssh
```

and optionally a specific user for netconf:

```text
set system login user netconf class xxxx
```
with authentication method : ssh key or password

```text
set system login user netconf authentication ssh-rsa "xxxx"
```
or

```text
set system login user netconf authentication plain-text-password
```

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Junos Provider
provider "junos" {
  ip         = var.junos_ip_or_dns
  sshkeyfile = var.ssh_key_path
}

# Configure an interface
resource junos_interface_physical "server1" {
  name = "ge-0/0/3"
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `ip` - (Required) This is the target for Netconf session (ip or dns name).  
  It can also be sourced from the `JUNOS_HOST` environment variable.

* `username` - (Optional) This is the username for ssh connection.  
  It can also be sourced from the `JUNOS_USERNAME` environment variable.  
  Defaults to `netconf`.

* `sshkey_pem` - (Optional) This is the ssh key in PEM format for establish ssh connection.  
  It can also be sourced from the `JUNOS_KEYPEM` environment variable.  
  Defaults is empty.

* `sshkeyfile` - (Optional) This is the path to ssh key for establish ssh connection.  
  Used only if `sshkey_pem` is empty.  
  It can also be sourced from the `JUNOS_KEYFILE` environment variable.  
  Defaults is empty.

* `password` - (Optional) This is a password for ssh connection.  
  Used only if `sshkey_pem` and `sshkeyfile` is empty.  
  It can also be sourced from the `JUNOS_PASSWORD` environment variable.  
  Defaults is empty.

* `port` - (Optional) This is the tcp port for ssh connection.  
  It can also be sourced from the `JUNOS_PORT` environment variable.  
  Defaults to `830`.

* `keypass` - (Optional) This is the passphrase for open `sshkeyfile` or `sshkey_pem`.  
  It can also be sourced from the `JUNOS_KEYPASS` environment variable.  
  Defaults is empty.

* `group_interface_delete` - (Optional) This is the Junos group used for remove configuration on a physical interface.  
  See interface specifications [interface specifications](#interface-specifications).  
  It can also be sourced from the `JUNOS_GROUP_INTERFACE_DELETE` environment variable.  
  Defaults to empty.

---
#### Command options
* `cmd_sleep_short` - (Optional) Number of milliseconds to wait after Terraform provider executes an action on the Junos device.  
  It can also be sourced from the `JUNOS_SLEEP_SHORT` environment variable.  
  Defaults to `100`.

* `cmd_sleep_lock` - (Optional) Number of seconds of standby while waiting for Terraform provider to lock candidate configuration on a Junos device.  
  It can also be sourced from the `JUNOS_SLEEP_LOCK` environment variable.  
  Defaults to `10`.

---
#### SSH options
* `ssh_sleep_closed` - (Optional) Number of seconds to wait after Terraform provider closed a ssh connection.  
  It can also be sourced from the `JUNOS_SLEEP_SSH_CLOSED` environment variable.  
  Defaults to `0`.

---
#### Debug & workaround options
* `file_permission` - (Optional) The permission to set for the created file (debug, setfile).  
  It can also be sourced from the `JUNOS_FILE_PERMISSION` environment variable.  
  Defaults to `0644`.

* `debug_netconf_log_path` - (Optional) more detailed log (netconf) in the specified file.  
  It can also be sourced from the `JUNOS_LOG_PATH` environment variable.  
  Defaults is empty.

* `fake_create_with_setfile` - (Optional, **don't use in normal terraform run**) When this option is set (with a path to a file), the normal process to create resources (netconf connection, precheck, generate/upload set lines in candidate configuration, commit, postcheck) skipped to generate set lines, append them to the specified file, and respond with a `fake` successful creation of resource to Terraform.  
Then you can upload/commit the file with the `junos_null_commit_file` resource in the same config or another terraform config or with another way.  
If you are using `junos_null_commit_file` in the same terraform config, you must create dependencies between resources so that the creation of the ` junos_null_commit_file` resource is alone and last.  
This options is useful to create a workaround for a long terraform run if there are many ressources to be created and Junos device is slow to commit.  
As many tests are skipped, this option may generate extra config (not managed by terraform) on Junos device or conficts/errors for resources in tfstate. A `terraform refresh` will be able to detect parts of errors but **be carefully with this option**.  
There are exceptions for ressources :
  * `junos_null_commit_file`, the skip doesn’t of course concern this resource.
  * `junos_interface_st0_unit` cannot take into account the option and run still normal process.
  * `junos_interface_physical` don’t generate `chassis aggregated-devices ethernet device-count` line when it should be necessary.

   It can also be sourced from the `JUNOS_FAKECREATE_SETFILE` environment variable.  
   Defaults is empty.

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

## Number of ssh connections and netconf commands

By default, terraform run with 10 parrallel actions, cf [walks the graph](https://www.terraform.io/docs/internals/graph.html#walking-the-graph).

With N for terraform's [`-parallelism`](https://www.terraform.io/docs/commands/plan.html#parallelism-n) argument, this provider :

* open N ssh connections.
* reduce the parrallelism of netconf `show` commands parrallelism under N with a mutex lock.
* lock the Junos configuration before adding `set` lines and execute `commit` so one `commit` at a time (other threads wait for locking).

To reduce :

* the rate of parallel ssh connections, reduce parallelism with terraform's [`-parallelism`](https://www.terraform.io/docs/commands/plan.html#parallelism-n) argument.
* the rate of new ssh connections by second, increase the provider's [`ssh_sleep_closed`](#ssh_sleep_closed) argument.
* the rate of netconf commands by second on ssh connections, increase the provider's [`cmd_sleep_short`](#cmd_sleep_short) argument.

To increase :

* the speed of `commit` (if your Junos device is quick to commit), decrease the provider's [`cmd_sleep_lock`](#cmd_sleep_lock) argument (be safe, too small is counterproductive).
