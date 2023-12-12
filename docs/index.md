---
page_title: "Provider: Junos"
---

# Junos Provider

The Junos provider communicate with Junos device via netconf protocol and modify a part of configuration.

The provider allows you to manage some elements on Junos device.

## Provider installation

For automatic installation (Terraform 0.13 and later) use [registry](https://registry.terraform.io/providers/jeremmfr/junos/):

```hcl
terraform {
  required_providers {
    junos = {
      source = "jeremmfr/junos"
    }
  }
}
```

For manual installation see [README on github](https://github.com/jeremmfr/terraform-provider-junos/#manual-install)

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
resource "junos_interface_physical" "server1" {
  name = "ge-0/0/3"
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

- **ip** (Required, String)  
  This is the target for Netconf session (ip or dns name).  
  It can also be sourced from the `JUNOS_HOST` environment variable.

- **username** (Optional, String)  
  This is the username for ssh connection.  
  It can also be sourced from the `JUNOS_USERNAME` environment variable.  
  Defaults to `netconf`.

- **sshkey_pem** (Optional, String)  
  This is the ssh key in PEM format for establish ssh connection.  
  It can also be sourced from the `JUNOS_KEYPEM` environment variable.  
  Defaults to empty.

- **sshkeyfile** (Optional, String)  
  This is the path to ssh key for establish ssh connection.  
  Used only if `sshkey_pem` is empty.  
  It can also be sourced from the `JUNOS_KEYFILE` environment variable.  
  Defaults to empty.

- **password** (Optional, String)  
  This is a password for ssh connection.  
  It can also be sourced from the `JUNOS_PASSWORD` environment variable.  
  Defaults to empty.

- **port** (Optional, Number)  
  This is the tcp port for ssh connection.  
  It can also be sourced from the `JUNOS_PORT` environment variable.  
  Defaults to `830`.

- **keypass** (Optional, String)  
  This is the passphrase for open `sshkeyfile` or `sshkey_pem`.  
  It can also be sourced from the `JUNOS_KEYPASS` environment variable.  
  Defaults to empty.

- **group_interface_delete** (Optional, String)  
  This is the Junos group used to remove configuration on a physical interface.  
  See interface specifications [interface specifications](#interface-specifications).  
  It can also be sourced from the `JUNOS_GROUP_INTERFACE_DELETE` environment variable.  
  Defaults to empty.

-> **Note:**
  Two SSH authentication methods (keys / password) are possible and tried with the `sshkey_pem`,
  `sshkeyfile` arguments or the keys provided by a SSH agent through the `SSH_AUTH_SOCK`
  environnement variable and `password` argument.  
  The keys provided by a SSH agent are only read if `sshkey_pem` and `sshkeyfile` arguments aren't set.

---

### Command options

- **cmd_sleep_short** (Optional, Number)  
  Milliseconds to wait after Terraform provider executes an action on the Junos device.  
  It can also be sourced from the `JUNOS_SLEEP_SHORT` environment variable.  
  Defaults to `100`.

- **cmd_sleep_lock** (Optional, Number)  
  Seconds of standby while waiting for Terraform provider to lock candidate configuration on a
  Junos device.  
  It can also be sourced from the `JUNOS_SLEEP_LOCK` environment variable.  
  Defaults to `10`.

- **commit_confirmed** (Optional, Number)  
  Number of minutes until automatic rollback (1..65535).  
  It can also be sourced from the `JUNOS_COMMIT_CONFIRMED` environment variable.  

  **If this argument is specified**, for each resource action with commit,
  the commit will take place in three steps:
  - commit with the `confirmed` option and with the value of this argument as `confirm-timeout`.
  - wait for `<commit_confirmed_wait_percent>`% of the minutes defined in the value of this argument.
  - confirm commit to avoid rollback with the `commit check` command.

  If a gracefully shutting down call with `Ctrl-c` is received by Terraform,
  the wait step is stopped and provider returns an error.

- **commit_confirmed_wait_percent** (Optional, Number)  
  Percentage of `<commit_confirmed>` minute(s) to wait between
  `commit confirmed` (commit with automatic rollback) and
  `commit check` (confirmation) commands (0..99).  
  No effect if `<commit_confirmed>` is not used.  
  It can also be sourced from the `JUNOS_COMMIT_CONFIRMED_WAIT_PERCENT` environment variable.  
  Defaults to `90`.

---

### SSH options

- **ssh_sleep_closed** (Optional, Number)  
  Seconds to wait after Terraform provider closed a ssh connection.  
  It can also be sourced from the `JUNOS_SLEEP_SSH_CLOSED` environment variable.  
  Defaults to `0`.

- **ssh_ciphers** (Optional, List of String)  
  Ciphers used in SSH connection.  
  Defaults to [
  `aes128-gcm@openssh.com`,
  `aes256-gcm@openssh.com`,
  `chacha20-poly1305@openssh.com`,
  `aes128-ctr`,
  `aes192-ctr`,
  `aes256-ctr`
  ]

- **ssh_timeout_to_establish** (Optional, Number)  
  Seconds to wait for establishing TCP connections when initiating SSH connections.  
  It can also be sourced from the `JUNOS_SSH_TIMEOUT_TO_ESTABLISH` environment variable.  
  Defaults to `0` (no timeout).

- **ssh_retry_to_establish** (Optional, Number)  
  Number of retries to establish SSH connections.  
  The provider waits after each try, with the sleep time increasing by 1 second each time.  
  It can also be sourced from the `JUNOS_SSH_RETRY_TO_ESTABLISH` environment variable.  
  Defaults to `1` (1..10).

---

### Debug & workaround options

- **file_permission** (Optional, String)  
  The permission to set for the created file (debug, setfile).  
  It can also be sourced from the `JUNOS_FILE_PERMISSION` environment variable.  
  Defaults to `0644`.

- **debug_netconf_log_path** (Optional, String)  
  More detailed log (netconf) in the specified file.  
  It can also be sourced from the `JUNOS_LOG_PATH` environment variable.  
  Defaults to empty.

  ~> **NOTE:** If this option is used (not empty), all Junos commands are logged in this file,
  therefore there may be sensitive data in plain text in the file.
  For example, when you use `plain_text_password` in the `junos_system_login_user` resource.

- **fake_create_with_setfile** (Optional, String, **don't use in normal terraform run**)
  When this option is set (with a path to a file), the normal process to create resources (netconf
  connection, pre-check, generate/upload set lines in candidate configuration, commit, post-check)
  skipped to generate set lines, append them to the specified file, and respond with a `fake`
  successful creation of resources to Terraform.  
  Then you can upload/commit the file with the `junos_null_commit_file` resource in the same config
  or another terraform config or with another way.  
  If you are using `junos_null_commit_file` in the same terraform config, you must create dependencies
  between resources so that the creation of the `junos_null_commit_file` resource is alone and
  last.  
  This options is useful to create a workaround for a long terraform run if there are many resources
  to be created and Junos device is slow to commit.  
  As many tests are skipped, this option may generate extra config (not managed by Terraform) on
  Junos device or conflicts/errors for resources in tfstate.
  A `terraform refresh` will be able to detect parts of errors but **be careful with**
  **this option**.  
  There are exceptions for resources :
  - **junos_interface_physical** don’t generate `chassis aggregated-devices ethernet device-count`
    line when it should be necessary.
  - **junos_interface_st0_unit** cannot take into account the option and run still
    normal process.
  - **junos_null_commit_file**, the skip doesn’t of course concern this resource.

  It can also be sourced from the `JUNOS_FAKECREATE_SETFILE` environment
  variable.  
  Defaults to empty.

  ~> **NOTE:** If this option is used (not empty), all Junos commands are added to this file,
  therefore there may be sensitive data in plain text in the file.
  For example, when you use `plain_text_password` in the `junos_system_login_user` resource.

- **fake_update_also** (Optional, Boolean, **don't use in normal terraform run**)  
  As with `create` and `fake_create_with_setfile`, when this option is true, the normal
  process to update resources skipped to generate set/delete lines, append them to the same file as
  `fake_create_with_setfile`, and respond with a `fake` successful update of resources to
  Terraform.  
  As with `fake_create_with_setfile`, this option may generate conflicts/errors for resources
  in tfstate. A `terraform refresh` will be able to detect parts of errors but
  **be careful with this option**.  
  There are exceptions for resources :
  - **junos_interface_physical** don’t generate `chassis aggregated-devices ethernet device-count`
    line when it should be necessary.
  - **junos_null_commit_file**, the skip doesn’t of course concern this resource.

  It can also be sourced from the `JUNOS_FAKEUPDATE_ALSO` environment variable and
  its value is `true`.  
  Defaults to `false`.

- **fake_delete_also** (Optional, Boolean, **don't use in normal terraform run**)  
  As with `create` and `fake_create_with_setfile`, when this option is true, the normal
  process to delete resources skipped to generate delete lines, append them to the same file as
  `fake_create_with_setfile`, and respond with a `fake` successful delete of resources to
  Terraform.  
  As with `fake_create_with_setfile`, this option may leave extra config (not managed by Terraform)
  on Junos device. **Be careful with this option**.  
  There are exceptions for resources :
  - **junos_interface_physical** don’t generate `chassis aggregated-devices ethernet device-count`
    line when it should be necessary.
  - **junos_null_commit_file**, the skip doesn’t of course concern this resource.

  It can also be sourced from the `JUNOS_FAKEDELETE_ALSO` environment variable and
  its value is `true`.  
  Defaults to `false`.

## Interface specifications

When create a resource for a physical interface, the provider considers the interface available if
there is ```apply-groups <group_interface_delete>``` and only this line on interface configuration.

Example if `group_interface_delete` = `interface-NC`:

```text
ge-0/0/3 {
  apply-groups interface-NC;
}
```

When provider destroy resource for physical interface, he add this line.

If `group_interface_delete` is empty the provider add this configuration on physical interface when
delete resource :

```text
ge-0/0/3 {
  description NC;
  disable;
}
```

and considers the interface available if there is this lines and only this lines on interface.

## Number of ssh connections and netconf commands

By default, terraform run with 10 parallel actions, cf [walks the graph](https://www.terraform.io/docs/internals/graph.html#walking-the-graph).

With N for Terraform's [`-parallelism`](https://www.terraform.io/docs/commands/plan.html#parallelism-n)
argument, this provider :

- open N ssh connections.
- reduce the parallelism of netconf `show` commands parallelism under N with a mutex lock.
- lock the Junos configuration before adding `set` lines and execute `commit` so one `commit` at a
time (other threads wait for locking).

To reduce :

- the rate of parallel ssh connections, reduce parallelism with Terraform's
[`-parallelism`](https://www.terraform.io/docs/commands/plan.html#parallelism-n) argument.
- the rate of new ssh connections by second, increase the provider's `ssh_sleep_closed` argument.
- the rate of netconf commands by second on ssh connections, increase the provider's
`cmd_sleep_short` argument.

To increase :

- the speed of `commit` (if your Junos device is quick to commit), decrease the provider's
`cmd_sleep_lock` argument (be safe, too small is counterproductive).
