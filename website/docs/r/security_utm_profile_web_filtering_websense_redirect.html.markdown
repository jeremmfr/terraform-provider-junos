---
layout: "junos"
page_title: "Junos: junos_security_utm_profile_web_filtering_websense_redirect"
sidebar_current: "docs-junos-resource-security-utm-profile-web-filtering-websense-redirect"
description: |-
  Create a security utm feature-profile web-filtering websense-redirect profile (when Junos device supports it)
---

# junos_security_utm_profile_web_filtering_websense_redirect

Provides a security utm feature-profile web-filtering websense-redirect profile resource.

## Example Usage

```hcl
# Add a security utm feature-profile web-filtering websense-redirect profile
resource junos_security_utm_profile_web_filtering_websense_redirect "demo_profile" {
  name                 = "Default Webfilter3"
  custom_block_message = "Blocked by Juniper"
  server {
    host = "10.0.0.1"
    port = 1024
  }
  timeout = 3
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource)(`String`) The name of security utm feature-profile web-filtering websense-redirect profile.
* `account` - (Optional)(`String`) Set websense redirect account.
* `custom_block_message` - (Optional)(`String`) Custom block message sent to HTTP client.
* `fallback_settings` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure fallback settings. Max of 1.
  * `default` - (Optional)(`String`) Default action. Need to be 'block' or 'log-and-permit'.
  * `server_connectivity` - (Optional)(`String`) Action when device cannot connect to server. Need to be 'block' or 'log-and-permit'.
  * `timeout` - (Optional)(`String`) Action when connection to server timeout. Need to be 'block' or 'log-and-permit'.
  * `too_many_requests` - (Optional)(`String`) Action when requests exceed the limit of engine. Need to be 'block' or 'log-and-permit'.
* `server` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure server settings. Max of 1.
  * `host` - (Optional)(`String`) Server host IP address or string host name.
  * `port` - (Optional)(`Int`) Server port. Need to be between 1024 and 65535.
* `socket` - (Optional)(`Int`) Set sockets number. Need to be between 1 and 32.
* `timeout` - (Optional)(`Int`) Set timeout. Need to be between 1 and 1800.

## Import

Junos security utm feature-profile web-filtering websense-redirect profile can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_profile_web_filtering_websense_redirect.demo_profile "Default Webfilter3"
```
