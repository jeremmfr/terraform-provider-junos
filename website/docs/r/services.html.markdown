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

* `application_identification` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'application-identification'. See the [`application_identification` arguments] (#application_identification-arguments) block.
* `security_intelligence` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'security-intelligence' configuration. See the [`security_intelligence` arguments] (#security_intelligence-arguments) block.

---
#### application_identification arguments
* `application_system_cache` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable application system cache. Conflict with `no_application_system_cache`.
  * `no_miscellaneous_services` - (Optional)(`Bool`) Disable ASC for miscellaneous services APBR,...
  * `security-services` - (Optional)(`Bool`) Enable ASC for security services (appfw, appqos, idp, skyatp..).
* `no_application_system_cache` - (Optional)(`Bool`) Disable storing AI result in application system cache. Conflict with `application_system_cache`.
* `application_system_cache_timeout` - (Optional)(`Int`) Application system cache entry lifetime (0..1000000).
* `download` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'download' configuration.
  * `automatic_interval` - (Optional)(`Int`) Attempt to download new application package (6..720 hours).
  * `automatic_start_time` - (Optional)(`String`) Start time to scheduled download and update (MM-DD.hh:mm / YYYY-MM-DD.hh:mm:ss).
  * `ignore_server_validation` - (Optional)(`Bool`) Disable server authentication for Applicaton Signature download.
  * `proxy_profile` - (Optional)(`String`) Configure web proxy for Application signature download
  * `url` - (Optional)(`String`) URL for application package download.
* `enable_performance_mode` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable performance mode knobs for best DPI performance.
  * `max_packet_threshold` - (Optional)(`Int`) Set the maximum packet threshold for DPI performance mode (1..100).
* `global_offload_byte_limit` - (Optional)(`Int`) Global byte limit to offload AppID inspection (0..4294967295).
* `imap_cache_size` - (Optional)(`Int`) IMAP cache size, it will be effective only after next appid sigpack install (60..512000).
* `imap_cache_timeout` - (Optional)(`Int`) IMAP cache entry timeout in seconds (1..86400).
* `inspection_limit_tcp` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable TCP byte/packet inspection limit.
  * `byte_limit` - (Optional)(`Int`) TCP byte inspection limit (0..4294967295).
  * `packet_limit` - (Optional)(`Int`) TCP packet inspection limit (0..4294967295).
* `inspection_limit_udp` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable UDP byte/packet inspection limit.
  * `byte_limit` - (Optional)(`Int`) UDP byte inspection limit (0..4294967295).
  * `packet_limit` - (Optional)(`Int`) UDP packet inspection limit (0..4294967295).
* `max_memory` - (Optional)(`Int`) Maximum amount of object cache memory JDPI can use (in MB) (1..200000).
* `max_transactions` - (Optional)(`Int`) Number of transaction finals to terminate application classification (0..25)
* `micro_apps` - (Optional)(`Bool`) Enable Micro Apps identifcation.
* `statistics_interval` - (Optional)(`Int`) Configure application statistics information with collection interval (1..1440 minutes).

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
