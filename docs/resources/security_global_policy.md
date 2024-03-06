---
page_title: "Junos: junos_security_global_policy"
---

# junos_security_global_policy

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `security policies global` block.

Configure static configuration in `security policies global` block

## Example Usage

```hcl
# Configure security policies global
resource "junos_security_global_policy" "global" {
  policy {
    name                      = "test"
    match_source_address      = ["blue"]
    match_destination_address = ["green"]
    match_application         = ["any"]
    match_from_zone           = ["any"]
    match_to_zone             = ["any"]
  }
  policy {
    name                      = "drop"
    match_source_address      = ["blue"]
    match_destination_address = ["any"]
    match_application         = ["any"]
    match_from_zone           = ["any"]
    match_to_zone             = ["any"]
    then                      = "deny"
  }
}
```

## Argument Reference

The following arguments are supported:

- **policy** (Required, Block List)  
  For each policy name.
  - **name** (Required, String)  
    Security policy name.
  - **match_source_address** (Required, Set of String)  
    List of source address match.
  - **match_destination_address** (Required, Set of String)  
    List of destination address match.
  - **match_from_zone** (Required, Set of String)  
    Match multiple source zone.
  - **match_to_zone** (Required, Set of String)  
    Match multiple destination zone.
  - **then** (Optional, String)  
    Action of policy.  
    Defaults to `permit`.
  - **count** (Optional, Boolean)  
    Enable count.
  - **log_init** (Optional, Boolean)  
    Log at session init time.
  - **log_close** (Optional, Boolean)  
    Log at session close time.
  - **match_application** (Optional, Set of String)  
    List of applications match.
  - **match_destination_address_excluded** (Optional, Boolean)  
    Exclude destination addresses.
  - **match_dynamic_application** (Optional, Set of String)  
    List of dynamic application or group match.
  - **match_source_address_excluded** (Optional, Boolean)  
    Exclude source addresses.
  - **match_source_end_user_profile** (Optional, String)  
    Match source end user profile (device identity profile).
  - **permit_application_services** (Optional, Block)  
    Declare `permit application-services` configuration.  
    See [below for nested schema](#permit_application_services-arguments).

---

### permit_application_services arguments

- **advanced_anti_malware_policy** (Optional, String)  
  Specify advanced-anti-malware policy name.
- **application_firewall_rule_set** (Optional, String)  
  Service rule-set name for Application firewall.
- **application_traffic_control_rule_set** (Optional, String)  
  Service rule-set name Application traffic control.
- **gprs_gtp_profile** (Optional, String)  
  Specify GPRS Tunneling Protocol profile name.
- **gprs_sctp_profile** (Optional, String)  
  Specify GPRS stream control protocol profile name.
- **idp** (Optional, Boolean)  
  Enable Intrusion detection and prevention.
- **idp_policy** (Optional, String)  
  Specify idp policy name.
- **redirect_wx** (Optional, Boolean)  
  Set WX redirection.
- **reverse_redirect_wx** (Optional, Boolean)  
  Set WX reverse redirection.
- **security_intelligence_policy** (Optional, String)  
  Specify security-intelligence policy name.
- **ssl_proxy** (Optional, Block)  
  Enable SSL Proxy.
  - **profile_name** (Optional, String)  
    Specify SSL proxy service profile name.
- **uac_policy** (Optional, Block)  
  Enable unified access control enforcement.
  - **captive_portal** (Optional, String)  
    Specify captive portal.
- **utm_policy** (Optional, String)  
  Specify utm policy name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `security_global_policy`.

## Import

Junos security global policies can be imported using any id, e.g.

```shell
$ terraform import junos_security_global_policy.global random
```
