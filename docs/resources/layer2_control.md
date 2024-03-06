---
page_title: "Junos: junos_layer2_control"
---

# junos_layer2_control

-> **Note:** This resource should only be created **once**.
It's used to configure options in `protocols layer2-control` block.  

Configure `protocols layer2-control` block

## Example Usage

```hcl
# Configure layer2 control
resource "junos_layer2_control" "l2control" {
  bpdu_block {}
}
```

## Argument Reference

The following arguments are supported:

- **bpdu_block** (Optional, Block)  
  Block BPDU on interface (BPDU Protect).  
  See [below for nested schema](#bpdu_block-arguments).
- **mac_rewrite_interface** (Optional, Block Set)  
  For each interface, Mac rewrite functionality.
  - **name** (Required, String)  
    Name of interface.
  - **enable_all_ifl** (Optional, Boolean)  
    Enable tunneling for all the IFLs under the interface.
  - **protocol** (Optional, Set of String)  
    Protocols for which mac rewrite need to be enabled.
- **nonstop_bridging** (Optional, Boolean)  
  Enable nonstop operation.

---

### bpdu_block arguments

- **disable_timeout** (Optional, Number)  
  Disable timeout for BPDU Protect (10..3600 seconds).
- **interface** (Optional, Block Set)  
  For each interface, to block BPDU on
  - **name** (Required, String)  
    Name of interface.
  - **disable** (Optional, Boolean)  
    Disable bpdu-block on a port.
  - **drop** (Optional, Boolean)  
    Drop xSTP BPDUs.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `layer2_control`.

## Import

Junos protocols layer2-control can be imported using any id, e.g.

```shell
$ terraform import junos_layer2_control.l2control random
```
