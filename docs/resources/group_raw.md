---
page_title: "Junos: junos_group_raw"
---

# junos_group_raw

Provides a group resource with raw configuration.

This resource allows you to create a Junos group with
raw configuration in either text or set format.  
The configuration is loaded directly without structured parsing,
giving you full flexibility to define any Junos configuration within a group.

## Example Usage

```hcl
# Create a group with text format (default)
resource "junos_group_raw" "dns_config" {
  name   = "dns-servers"
  config = <<EOT
system {
    services {
        dns {
            forwarders {
                192.0.2.3;
                192.0.2.33;
            }
        }
    }
}
EOT
}

# Create a group with set format
resource "junos_group_raw" "system_config" {
  name   = "system-settings"
  format = "set"
  config = <<EOT
set system time-zone Europe/Paris
set system default-address-selection
set system ntp peer 192.0.2.1
EOT
}

# Apply the group
resource "junos_apply_group" "dns" {
  name = junos_group_raw.dns_config.name
}

resource "junos_apply_group" "system" {
  name   = junos_group_raw.system_config.name
  prefix = "system "
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of the group.
- **config** (Required, String)  
  The raw configuration to load.  
  The format of this configuration depends on the `format` attribute.
- **format** (Optional, String, Forces new resource)  
  The format used for the configuration data.  
  Need to be `text` or `set`.  
  Defaults to `text`.
  - When `text`: configuration should be in Junos text format (curly braces).
  - When `set`: each line must start with `set ` and be in Junos set command format.

## Attribute Reference

The following attributes are exported:

- **id** (String)
  An identifier for the resource with format `<name>_-_<format>`.

## Import

Junos group can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_group_raw.dns_config "dns-servers"
```

Or with the format specified `<name>_-_<format>`:

```shell
$ terraform import junos_group_raw.system_config "system-settings_-_set"
```
