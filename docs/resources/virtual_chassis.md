---
page_title: "Junos: junos_virtual_chassis"
---

# junos_virtual_chassis

-> **Note:** This resource should only be created **once**.
It's used to configure `virtual-chassis` block.  

Configure `virtual-chassis` block

## Example Usage

```hcl
# Configure virtual-chassis
resource "junos_virtual_chassis" "virtual_chassis" {
  no_split_detection = true
}
```

## Argument Reference

The following arguments are supported:

- **alias** (Optional, Block Set)  
  For each serial_number, provide an alias name for this serial-number.
  - **serial_number** (Required, String)  
    Member's serial number.
  - **alias_name** (Required, String)  
    Alias name for this serial-number.
- **auto_sw_update** (Optional, Boolean)  
  Auto software update.
- **auto_sw_update_package_name** (Optional, String)  
  URL or pathname of software package to auto software update.  
  `auto_sw_update` need to be true.
- **graceful_restart_disable** (Optional, Boolean)  
  Disable graceful restart.
- **identifier** (Optional, String)  
  Virtual chassis identifier, of type ISO system-id.
- **mac_persistence_timer** (Optional, String)  
  MAC persistence time (minutes) or disable.  
  Need to be a number between 1 to 60 or `disable`.
- **member** (Optional, Block List)  
  For each identifier, member of virtual chassis configuration.  
  See [below for nested schema](#member-arguments).
- **no_split_detection** (Optional, Boolean)  
  Disable split detection.
- **preprovisioned** (Optional, Boolean)  
  Only accept preprovisioned members.
- **traceoptions** (Optional, Block)  
  Trace options for virtual chassis.  
  See [below for nested schema](#traceoptions-arguments).
- **vcp_no_hold_time** (Optional, Boolean)  
  Set no hold time for vcp interfaces.

### member arguments

- **id** (Required, Number)  
  Member identifier (0..9).
- **location** (Optional, String)  
  Member's location.
- **mastership_priority** (Optional, Number)  
  Member's mastership priority.
- **no_management_vlan** (Optional, Boolean)  
  Disable management VLAN.
- **role** (Optional, String)  
  Member's role.  
  Need to be `line-card` or `routing-engine`.
- **serial_number** (Optional, String)  
  Member's serial number.

### traceoptions arguments

- **flag** (Optional, Set of String)  
  Tracing parameters.
- **file** (Optional, Block)  
  Declare `file` configuration.
  - **name** (Required, String)  
    Name of file in which to write trace information.
  - **files** (Optional, Number)  
    Maximum number of trace files (2..1000).
  - **no_stamp** (Optional, Boolean)  
    Do not timestamp trace file.
  - **replace** (Optional, Boolean)  
    Replace trace file rather than appending to it.
  - **size** (Optional, Number)  
    Maximum trace file size.
  - **world_readable** (Optional, Boolean)  
    Allow any user to read the log file.
  - **no_world_readable** (Optional, Boolean)  
    Don't allow any user to read the log file.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `virtual-chassis`.

## Import

Junos virtual-chassis can be imported using any id, e.g.

```shell
$ terraform import junos_virtual_chassis.virtual_chassis random
```
