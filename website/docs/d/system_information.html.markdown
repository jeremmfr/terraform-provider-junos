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

* `id` - Hostname of the Junos device
* `hardware_model` - Type of hardware/software of Junos device (i.e. - SRX340, vSRX, etc)
* `os_name` - Operating system name of Junos
* `os_version` - Software version of Junos
* `serial_number` - Serial number of the device
* `cluster_node` - Boolean flag that indicates if device is part of a cluster or not
