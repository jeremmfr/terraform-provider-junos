---
page_title: "Junos: junos_policyoptions_community"
---

# junos_policyoptions_community

Get configuration from a policy-options community.

## Example Usage

```hcl
# Read a policy-options community configuration
data "junos_policyoptions_community" "community_demo" {
  name = "communityDemo"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String)  
  Name to identify BGP community.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the data source with format `<name>`.
- **dynamic_db** (Boolean)  
  Object may exist in dynamic database.
- **invert_match** (Boolean)  
  Invert the result of the community expression matching.
- **members** (List of String)  
  Community members.
