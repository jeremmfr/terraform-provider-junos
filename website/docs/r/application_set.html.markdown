---
layout: "junos"
page_title: "Junos: junos_application_set"
sidebar_current: "docs-junos-resource-application-set"
description: |-
  Create a set of application
---

# junos_application_set

Provides a application set.

## Example Usage

```hcl
# Add a set of application
resource junos_application_set "ssh_telnet" {
  name         = "ssh_telnet"
  applications = ["junos-ssh", "junos-telnet"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Name of application set.
* `applications` - (Optional)(`ListOfString`) List of application names.

## Import

Junos application set can be imported using an id made up of `<name>`, e.g.

```
$ terraform import junos_application_set.ssh_telnet ssh_telnet
```
