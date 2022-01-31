---
page_title: "Junos: junos_policyoptions_community"
---

# junos_policyoptions_community

Provides a community BGP resource.

## Example Usage

```hcl
# Add a community
resource junos_policyoptions_community "community_demo" {
  name    = "communityDemo"
  members = ["65000:100"]
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of community.
- **members** (Required, List of String)  
  List of community.
- **invert_match** (Optional, Boolean)  
  Add `invert-match` parameter.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos community can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_policyoptions_community.community_demo communityDemo
```
