---
page_title: "Junos: junos_services_security_intelligence_profile"
---

# junos_services_security_intelligence_profile

Provides a services security-intelligence profile resource.

## Example Usage

```hcl
# Add a services security-intelligence profile
resource "junos_services_security_intelligence_profile" "demo" {
  name     = "demo"
  category = "CC"
  rule {
    name = "rule_1"
    match {
      threat_level = [10]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Security intelligence profile name.
- **category** (Required, String)  
  Profile category name.
- **rule** (Required, Block List)  
  For each rule name.
  See [below for nested schema](#rule-arguments).
- **default_rule_then** (Optional, Block)  
  Declare profile default rule.
  - **action** (Required, String)  
    Security intelligence profile action.  
    Need to be `permit`, `recommended`, `block drop`, `block close` or
    `block close http (file|message|redirect-url) ...`.
  - **log** (Optional, Boolean)  
    Log security intelligence block action.
  - **no_log** (Optional, Boolean)  
    Don't log security intelligence block action.
- **description** (Optional, String)  
  Text description of profile.

---

### rule arguments

- **name** (Required, String)  
  Profile rule name.
- **match** (Required, Block)  
  Configure profile matching feed name and threat levels.
  - **threat_level** (Required, List of Number)  
    Profile matching threat levels, higher number is more severe (1..10).
  - **feed_name** (Optional, List of String)  
    Profile matching feed name.
- **then_action** (Required, String)  
  Security intelligence profile action.  
  Need to be `permit`, `recommended`, `block drop`, `block close` or
  `block close http (file|message|redirect-url) ...`.
- **then_log** (Optional, Boolean)  
  Log security intelligence block action.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos services security-intelligence profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_services_security_intelligence_profile.demo demo
```
