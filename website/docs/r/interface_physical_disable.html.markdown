---
layout: "junos"
page_title: "Junos: junos_interface_physical_disable"
sidebar_current: "docs-junos-resource-interface-physical-disable"
description: |-
  Disable a not configured physical interface (same as when destroy junos_interface_physical resource)
---

# junos_interface_physical_disable

Disable a not configured physical interface
(same as when destroy `junos_interface_physical` resource).  
If the interface is configured or is used for a logical unit interface, the apply fails.

This resource is useful for disable physical interfaces that have not already been used once
by the `junos_interface_physical` resource.

Destroy this resource has no effect on the Junos configuration.

## Example Usage

```hcl
# Disable an interface
resource junos_interface_physical_disable "interface_demo" {
  name = "ge-0/0/0"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of physical interface (without dot).

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.
