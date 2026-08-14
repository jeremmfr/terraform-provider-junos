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

```hcl
# Configure chassis FPC slot 0 with the sampling instance attached by the sampling instance resource
resource "junos_chassis_fpc" "fpc0" {
  slot_number                        = 0
  configure_sampling_instance_singly = true

  error {
    major_action    = "alarm"
    major_threshold = 10
  }
}

resource "junos_forwardingoptions_sampling_instance" "demo" {
  name                     = "demo"
  chassis_fpc_slot_numbers = [0]

  input {
    rate = 1
  }
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
- **configure_sampling_instance_singly** (Optional, Boolean)  
  Configure `sampling-instance` option in other resource
  (like `junos_forwardingoptions_sampling_instance` with `chassis_fpc_slot_numbers` argument).  
  Conflict with `sampling_instance`.  
  Required on some Junos systems where the sampling instance and the FPC binding cannot be
  committed separately.
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
