---
layout: "junos"
page_title: "Junos: junos_eventoptions_destination"
sidebar_current: "docs-junos-resource-eventoptions-destination"
description: |-
  Create an event-options destination
---

# junos_eventoptions_destination

Provides an event-options destination resource.

## Example Usage

```hcl
# Add an event-options destination
resource "junos_eventoptions_destination" "demo" {
  name = "demo"
  archive_site {
    url = "https://example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) Destination name.
* `archive_site` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified multiple times for each archive destination.
  * `url` - (Required)(`String`) URL of destination for file.
  * `password` - (Optional)(`String`) Password for login into the archive site.  
  **WARNING** Clear in tfstate.
* `transfer_delay` - (Optional)(`Int`) Delay before transferring files (seconds).

## Import

Junos event-options destination can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_destination.demo demo
```
