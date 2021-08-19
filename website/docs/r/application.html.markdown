---
layout: "junos"
page_title: "Junos: junos_application"
sidebar_current: "docs-junos-resource-application"
description: |-
  Create a application
---

# junos_application

Provides a application resource.

## Example Usage

```hcl
# Add a application
resource junos_application "mysql" {
  name             = "mysql"
  protocol         = "tcp"
  destination_port = "3306"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of application.
- **application_protocol** (Optional, String)  
  Application protocol type.
- **description** (Optional, String)  
  Text description of application.
- **destination_port** (Optional, String)  
  Port(s) destination used by application.
- **ether_type** (Optional, String)  
  Match ether type.  
  Must be in hex (example: 0x8906).
- **inactivity_timeout** (Optional, Number)  
  Application-specific inactivity timeout (4..86400 seconds).
- **protocol** (Optional, String)  
  Protocol used by application.
- **rpc_program_number** (Optional, String)  
  Match range of RPC program numbers.  
  Must be an integer or a range of integers.
- **source_port** (Optional, String)  
  Port(s) source used by application.
- **uuid** (Optional, String)  
  Match universal unique identifier for DCE RPC objects.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos application can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_application.mysql mysql
```
