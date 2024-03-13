---
page_title: "Junos: junos_routes"
---

# junos_routes

Get all routes in all routing tables or a selected routing table.

## Example Usage

```hcl
# Read all routes 
data "junos_routes" "all" {}
output "all_routes" {
  value = data.junos_routes.all.table
}
```

## Argument Reference

The following arguments are supported:

- **table_name** (Optional, String)  
  Get routes only on a specific routing table with the name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<table_name>` or `all` if not set.
- **table** (Block List)  
  For each routing table.
  - **name** (String)  
    Name of the routing table.
  - **route** (Block List)  
    For each route destination in the routing table.  
    See [below for nested schema](#route-attributes).

### route attributes

- **destination** (String)  
  Route destination.
- **entry** (Block List)  
  For each route to the destination.
  - **as_path** (String)  
    AS path through which the route was learned.
  - **current_active** (Boolean)  
    This entry is the current active one.
  - **local_preference** (Number)  
    Local preference value for the route.
  - **metric** (Number)  
    Cost value for the route.
  - **next_hop** (Block List)  
    For each next hop.  
    See [below for nested schema](#next_hop-attributes).
  - **next_hop_type** (String)  
    Next hop type.
  - **preference** (Number)  
    Preference value for the route.
  - **protocol** (String)  
    Protocol from which the route was learned.

### next_hop attributes

- **local_interface** (String)  
  Interface used to local routes.
- **selected_next_hop** (Boolean)  
  It's the currently used next-hop.
- **to** (String)  
  Next hop to the destination.
- **via** (String)  
  Interface used to reach the next hop.
