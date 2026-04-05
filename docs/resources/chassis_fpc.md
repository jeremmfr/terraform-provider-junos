---
page_title: "Junos: junos_chassis_fpc"
---

# junos_chassis_fpc

Provides a chassis FPC resource.

-> **Note**
  Unlike most resources, this resource only manages the attributes defined in its schema.  
  Any other configuration present under the same `chassis fpc <slot_number>` block is left untouched.

## Example Usage

```hcl
# Configure chassis FPC slot 0
resource "junos_chassis_fpc" "fpc0" {
  slot_number       = 0
  sampling_instance = junos_forwardingoptions_sampling_instance.demo.name
}
```

## Argument Reference

-> **Note**
  At least one of arguments need to be set (in addition to `slot_number`).

The following arguments are supported:

- **slot_number** (Required, Number, Forces new resource)  
  FPC number.
- **cfp_to_et** (Optional, Boolean)  
  Enable ET interface and remove CFP client.
- **sampling_instance** (Optional, String)  
  Name for sampling instance.
- **error** (Optional, Block)  
  Error level configuration for FPC.  
  See [below for nested schema](#error-arguments).

---

### error arguments

- **fatal_action** (Optional, String)  
  Configure the action for fatal level.  
  Need to be `alarm`, `disable-pfe`, `get-state`, `log`, `offline`, `reset` or `trap`.
- **fatal_threshold** (Optional, Number)  
  Error count at which to take the action (1..1024).
- **major_action** (Optional, String)  
  Configure the action for major level.  
  Need to be `alarm`, `disable-pfe`, `get-state`, `log`, `offline`, `reset` or `trap`.
- **major_threshold** (Optional, Number)  
  Error count at which to take the action (1..1024).
- **minor_action** (Optional, String)  
  Configure the action for minor level.  
  Need to be `alarm`, `disable-pfe`, `get-state`, `log`, `offline`, `reset` or `trap`.
- **minor_threshold** (Optional, Number)  
  Error count at which to take the action (0..1024).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<slot_number>`.

## Import

Junos chassis FPC can be imported using the slot number as id, e.g.

```shell
$ terraform import junos_chassis_fpc.fpc0 0
```
