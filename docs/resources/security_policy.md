---
page_title: "Junos: junos_security_policy"
---

# junos_security_policy

Provides a security policy resource.

## Example Usage

```hcl
# Add a security policy
resource "junos_security_policy" "demo_policy" {
  from_zone = "trust"
  to_zone   = "untrust"
  policy {
    name                      = "allow_trust"
    match_source_address      = ["any"]
    match_destination_address = ["any"]
    match_application         = ["any"]
  }
}
```

## Argument Reference

The following arguments are supported:

- **from_zone** (Required, String, Forces new resource)  
  The name of source zone.
- **to_zone** (Required, String, Forces new resource)  
  The name of destination zone.
- **policy** (Required, Block List)  
  For each name of policy.
  - **name** (Required, String)  
    The name of policy.
  - **match_source_address** (Required, Set of String)  
    List of source address match.
  - **match_destination_address** (Required, Set of String)  
    List of destination address match.
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
  - **permit_tunnel_ipsec_vpn** (Optional, String)  
    Name of vpn to permit with a tunnel ipsec.
  - **permit_application_services** (Optional, Block)  
    Define application services for permit.  
    See [below for nested schema](#permit_application_services-arguments-for-policy).

---

### permit_application_services arguments for policy

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
  An identifier for the resource with format `<from_zone>_-_<to_zone>`.

## Import

Junos security policy can be imported using an id made up of `<from_zone>_-_<to_zone>`, e.g.

```shell
$ terraform import junos_security_zone.demo_policy trust_-_untrust
```
