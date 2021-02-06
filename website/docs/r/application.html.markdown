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

* `name` - (Required, Forces new resource)(`String`) Name of application.
* `destination_port` - (Optional)(`String`) Port(s) destination used by application.
* `protocol` - (Optional)(`String`) Protocol used by application.
* `source_port` - (Optional)(`String`) Port(s) source used by application.

## Import

Junos application can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_application.mysql mysql
```
