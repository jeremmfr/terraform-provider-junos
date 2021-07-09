---
layout: "junos"
page_title: "Junos: junos_switch_options"
sidebar_current: "docs-junos-resource-switch-options"
description: |-
  Configure static configuration in switch-options block
---

# junos_switch_options

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `switch-options` block. By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `switch-options` block

## Example Usage

```hcl
# Configure switch-options
resource junos_switch_options "switch_options" {
  vtep_source_interface = "lo0.0"
}
```

## Argument Reference

The following arguments are supported:

* `clean_on_destroy` - (Optional)(`Bool`) Clean supported lines when destroy this resource.
* `vtep_source_interface` - (Optional)(`String`) Source layer-3 IFL for VXLAN.

## Import

Junos switch_options can be imported using any id, e.g.

```shell
$ terraform import junos_switch_options.switch_options random
```
