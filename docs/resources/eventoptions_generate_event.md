---
page_title: "Junos: junos_eventoptions_generate_event"
---

# junos_eventoptions_generate_event

Provides an event-options generate-event resource.

## Example Usage

```hcl
# Add an event-options generate-event
resource "junos_eventoptions_generate_event" "demo" {
  name        = "demo"
  time_of_day = "12:00:00"
}
```

## Argument Reference

-> **Note:** One of `time_interval` or `time_of_day` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of the event to be generated.
- **no_drift** (Optional, Boolean)  
  Avoid event generation delay propagating to next event.
- **start_time** (Optional, String)  
  Start-time to generate event (YYYY-MM-DD.HH:MM:SS).  
  `time_interval` need to be set.
- **time_interval** (Optional, Number)  
  Frequency for generating the event (60..2592000 seconds).
- **time_of_day** (Optional, String)  
  Time of day at which to generate event (HH:MM:SS).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos event-options generate-event can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_generate_event.demo demo
```
