---
layout: "junos"
page_title: "Junos: junos_services"
sidebar_current: "docs-junos-resource-services"
description: |-
  Configure static configuration in services block
---

# junos_services

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `services` block. Destroy this resource has no effect on the Junos configuration.

Configure static configuration in `services` block

## Example Usage

```hcl
# Configure services
resource junos_services "services" {
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123400"
    url                  = "https://example.com/api/manifest.xml"
  }
}
```

## Argument Reference

The following arguments are supported:

* `security_intelligence` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'security-intelligence' configuration. See the [`security_intelligence` arguments] (#security_intelligence-arguments) block.

---
#### security_intelligence arguments
* `authentication_token` - (Optional)(`String`) Token string for authentication to use feed update services. Conflict with `authentication_tls_profile`.
* `authentication_tls_profile` - (Optional)(`String`) TLS profile for authentication to use feed update services. Conflict with `authentication_token`.
* `category_disable` - (Optional)(`ListOfString`) Categories to be disabled
* `default_policy` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure default-policy for a category. Can be specified multiple times for each category.
  * `category_name` - (Optional)(`String`) Name of security intelligence category.
  * `profile_name` - (Optional)(`String`) Name of profile.
* `proxy_profile` - (Optional)(`String`) The proxy profile name.
* `url` - (Optional)(`String`) Configure the url of feed server [https://<ip or hostname>:<port>/<uri>].
* `url_parameter` - (Optional)(`String`) Configure the parameter of url.

## Import

Junos services can be imported using any id, e.g.

```
$ terraform import junos_services.services random
```
