---
page_title: "Junos: junos_apply_group"
---

# junos_apply_group

Provides a resource to apply a Junos group at a specific configuration level.

This resource allows you to apply a group at the global level or
at a specific configuration prefix path.

-> **Note**
  There is no verification that the `prefix` path exists in the Junos configuration.
  Make sure the prefix is valid before applying the group.

## Example Usage

```hcl
# Create a group with raw configuration
resource "junos_group_raw" "dns_config" {
  name   = "dns-servers"
  config = <<EOT
system {
    services {
        dns {
            forwarders {
                192.0.2.3;
                192.0.2.4;
            }
        }
    }
}
EOT
}

# Apply the group globally
resource "junos_apply_group" "dns_global" {
  name = junos_group_raw.dns_config.name
}

# Apply the group at system level
resource "junos_apply_group" "dns_system" {
  name   = junos_group_raw.dns_config.name
  prefix = "system "
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of the group to apply.
- **prefix** (Optional, String, Forces new resource)  
  Prefix path to define where apply-group must be set.  
  If not specified, the group is applied globally.
  The prefix must end with a space character.

## Attribute Reference

The following attributes are exported:

- **id** (String)
  An identifier for the resource with format `<name>_-_<prefix>`.

## Import

Junos apply-group can be imported using an id made up of `<name>_-_<prefix>`, e.g.

```shell
$ terraform import junos_apply_group.dns_global "dns-servers_-_"
```

For a group applied with a prefix:

```shell
$ terraform import junos_apply_group.dns_system "dns-servers_-_system "
```
