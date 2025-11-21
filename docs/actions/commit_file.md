---
page_title: "Junos: junos_commit_file"
---

# junos_commit_file

Load a file with set/delete lines on device and commit.

This action provides a way to load and commit configuration from a file
containing Junos set/delete commands without creating a persistent resource in the Terraform state.

<!-- markdownlint-disable -->
-> **Note**
  Actions are a Terraform 1.14+ feature that allow you to perform operations without managing state.
  Unlike the `junos_null_commit_file` resource, this action does not create any state entries.
<!-- markdownlint-restore -->

## Example Usage

```hcl
action "junos_commit_file" "setfile" {
  config {
    filename                = "~/junos/setfile"
    clear_file_after_commit = true
  }
}
```

## Argument Reference

The following arguments are supported:

- **filename** (Required, String)  
  The path of the file to load.  
  Tilde (~) in the path will be expanded to the user's home directory.
- **append_lines** (Optional, List of String)  
  List of lines to append to the lines in the loaded file.
- **clear_file_after_commit** (Optional, Boolean)  
  Truncate file after successful commit.

## Progress Events

This action sends progress updates during execution:

- Reading configuration file
- Starting session to device
- Locking candidate configuration
- Loading configuration
- Committing configuration
- Configuration loaded and committed
- Clearing file after commit (if `clear_file_after_commit` is enabled)

## File Format

The file should contain Junos set and/or delete commands, one per line:

```text
set system host-name vSRX-1
set system time-zone Europe/Paris
delete system ntp server 192.0.2.1
set system ntp server 192.0.2.111
```
