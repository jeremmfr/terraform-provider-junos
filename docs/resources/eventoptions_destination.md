---
page_title: "Junos: junos_eventoptions_destination"
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

- **name** (Required, String, Forces new resource)  
  Destination name.
- **archive_site** (Required, Block List)  
  For each archive destination.
  - **url** (Required, String)  
    URL of destination for file.
  - **password** (Optional, String, Sensitive)  
    Password for login into the archive site.  
    Conflict with `password_wo`.
  - **password_wo** (Optional, String, Sensitive, Write-only)  
    Password for login into the archive site, not stored in state.  
    Requires `password_wo_version` and Terraform 1.11 or later.  
    Conflict with `password`.
  - **password_wo_version** (Optional, Number)  
    Version of `password_wo` to trigger the sending of its value.  
    Increment it to send the current value of `password_wo` to the device.  
    Requires `password_wo`.
- **transfer_delay** (Optional, Number)  
  Delay before transferring files (seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos event-options destination can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_destination.demo demo
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the passwords read on the device are
  stored in `password`, and therefore in the Terraform state.  
  When the configuration uses the write-only arguments, the next apply removes them from the state.
