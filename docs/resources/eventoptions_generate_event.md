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

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of the event to be generated.
- **time_interval** (Optional, Number)  
  Frequency for generating the event (60..2592000 seconds).  
  Need to set one of `time_interval` or `time_of_day`.
- **time_of_day** (Optional, String)  
  Time of day at which to generate event (HH:MM:SS).  
  Need to set one of `time_interval` or `time_of_day`.
- **no_drift** (Optional, Boolean)  
  Avoid event generation delay propagating to next event

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos event-options generate-event can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_eventoptions_generate_event.demo demo
```
