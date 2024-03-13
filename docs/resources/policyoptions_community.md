---
page_title: "Junos: junos_policyoptions_community"
---

# junos_policyoptions_community

Provides a policy-options community resource.

## Example Usage

```hcl
# Add a policy-options community
resource "junos_policyoptions_community" "community_demo" {
  name    = "communityDemo"
  members = ["65000:100"]
}
```

## Argument Reference

-> **Note:** One of `dynamic_db` or `members` arguments is required.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name to identify BGP community.
- **dynamic_db** (Optional, Boolean)  
  Object may exist in dynamic database.
- **invert_match** (Optional, Boolean)  
  Invert the result of the community expression matching.
- **members** (Optional, List of String)  
  Community members.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos community can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_community.community_demo communityDemo
```
