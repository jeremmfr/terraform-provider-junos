---
layout: "junos"
page_title: "Junos: system_information"
sidebar_current: "docs-junos-data-source-system-information"
description: |-
  Get information of the Junos device system information
---

# junos_system_information

Get information on the junos system

## Example Usage

```hcl
data junos_system_information "example" {}
```

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  Hostname of the Junos device.
- **hardware_model** (String)  
  Type of hardware/software of Junos device (i.e. - SRX340, vSRX, etc).
- **os_name** (String)  
  Operating system name of Junos.
- **os_version** (String)  
  Software version of Junos.
- **serial_number** (String)  
  Serial number of the device.
- **cluster_node** (Boolean)  
  Boolean flag that indicates if device is part of a cluster or not.
