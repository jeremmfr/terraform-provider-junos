---
page_title: "Junos: junos_apply_group_except"
---

# junos_apply_group_except

Provides a resource to exclude a Junos group at a specific configuration level.

This resource allows you to exclude a group from being applied at a specific
configuration prefix path using the `apply-groups-except` statement.

-> **Note**
  Unlike `junos_apply_group`, `apply-groups-except` cannot be applied globally.
  A `prefix` is always required to specify where the group exclusion should be set.

-> **Note**
  There is no verification that the `prefix` path exists in the Junos configuration.
  Make sure the prefix is valid before applying the group exclusion.

## Example Usage

```hcl
# Create a group with raw configuration
resource "junos_group_raw" "system" {
  name   = "system-default"
  format = "set"
  config = <<EOT
set system time-zone Europe/Paris
set system default-address-selection
EOT
}

# Apply the group globally
resource "junos_apply_group" "base" {
  name = junos_group_raw.system.name
}

# Exclude the group from system level
resource "junos_apply_group_except" "base_system" {
  name   = junos_group_raw.system.name
  prefix = "system "
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of the group to exclude.
- **prefix** (Required, String, Forces new resource)  
  Prefix path to define where apply-groups-except must be set.  
  The prefix must end with a space character.

## Attribute Reference

The following attributes are exported:

- **id** (String)
  An identifier for the resource with format `<name>_-_<prefix>`.

## Import

Junos apply-groups-except can be imported using an id made up of `<name>_-_<prefix>`, e.g.

```shell
$ terraform import junos_apply_group_except.base_system "system-default_-_system "
```
