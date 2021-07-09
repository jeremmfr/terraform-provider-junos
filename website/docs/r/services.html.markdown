---
layout: "junos"
page_title: "Junos: junos_services"
sidebar_current: "docs-junos-resource-services"
description: |-
  Configure static configuration in services block
---

# junos_services

-> **Note:** This resource should only be created **once**. It's used to configure static (not object) options in `services` block. By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

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

* `clean_on_destroy` - (Optional)(`Bool`) Clean supported lines when destroy this resource.
* `advanced_anti_malware` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'advanced-anti-malware' static configuration. See the [`advanced_anti_malware` arguments] (#advanced_anti_malware-arguments) block.
* `application_identification` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'application-identification'. See the [`application_identification` arguments] (#application_identification-arguments) block.
* `security_intelligence` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'security-intelligence' configuration. See the [`security_intelligence` arguments] (#security_intelligence-arguments) block.
* `user_identification` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'user-identification' configuration. See the [`user_identification` arguments] (#user_identification-arguments) block.

---

### advanced_anti_malware arguments

* `connection` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'connection' configuration.
  * `auth_tls_profile` - (Optional)(`String`) Authentication TLS profile.  
  **Note:** If not set, tls-profile is only read from the Junos configuration (so as not to be in conflict with enrollment process).
  * `proxy_profile` - (Optional)(`String`) Proxy profile.
  * `source_address` - (Optional)(`String`) The source ip for connecting to the cloud server. Conflict with `source_interface`.
  * `source_interface` - (Optional)(`String`) The source interface for connecting to the cloud server. Conflict with `source_address`.
  * `url` - (Optional)(`String`) The url of the cloud server [https://`<ip or hostname>`:`<port>`].  
  **Note:** If not set, url is only read from the Junos configuration (so as not to be in conflict with enrollment process).
* `default_policy` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'default-policy' configuration.
  * `blacklist_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware blacklist hit.
  * `default_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware action.
  * `fallback_options_action` - (Optional)(`String`) Notification action taken for fallback action. Need to be 'block' or 'permit'.
  * `fallback_options_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware fallback action.
  * `http_action` - (Optional)(`String`) Action taken for contents with verdict meet threshold for HTTP. Need to be 'block' or 'permit'. Required with `http_inspection_profile`.
  * `http_client_notify_file` - (Optional)(`String`) File name for http response to client notification action taken for contents with verdict meet threshold. Conflict with others `http_client_notify_*`. Required with `http_action` and `http_inspection_profile`.
  * `http_client_notify_message` - (Optional)(`String`) Block message to client notification action taken for contents with verdict meet threshold. Conflict with others `http_client_notify_*`. Required with `http_action` and `http_inspection_profile`.
  * `http_client_notify_redirect_url` - (Optional)(`String`) Redirect url to client notification action taken for contents with verdict meet threshold. Conflict with others `http_client_notify_*`. Required with `http_action` and `http_inspection_profile`.
  * `http_file_verdict_unknown` - (Optional)(`String`) Action taken for contents with verdict unknown. Required with `http_action` and `http_inspection_profile`.
  * `http_inspection_profile` - (Optional)(`String`) Advanced Anti-malware inspection-profile name for HTTP. Required with `http_action`.
  * `http_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware actions for HTTP. Required with `http_action` and `http_inspection_profile`.
  * `imap_inspection_profile` - (Optional)(`String`) Advanced Anti-malware inspection-profile name for IMAP.
  * `imap_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware actions for IMAP. Required with `imap_inspection_profile`.
  * `smtp_inspection_profile` - (Optional)(`String`) Advanced Anti-malware inspection-profile name for SMTP.
  * `smtp_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware actions for SMTP. Required with `smtp_inspection_profile`.
  * `verdict_threshold` - (Optional)(`String`) Verdict threshold.
  * `whitelist_notification_log` - (Optional)(`Bool`) Logging option for Advanced Anti-malware whitelist hit.

---

### application_identification arguments

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

### security_intelligence arguments

* `authentication_token` - (Optional)(`String`) Token string for authentication to use feed update services. Conflict with `authentication_tls_profile`.  
  **Note:** If not set, token is only read from the Junos configuration (so as not to be in conflict with enrollment process).
* `authentication_tls_profile` - (Optional)(`String`) TLS profile for authentication to use feed update services. Conflict with `authentication_token`.  
  **Note:** If not set, tls-profile is only read from the Junos configuration (so as not to be in conflict with enrollment process).
* `category_disable` - (Optional)(`ListOfString`) Categories to be disabled
* `default_policy` - (Optional)[attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Configure default-policy for a category. Can be specified multiple times for each category.
  * `category_name` - (Optional)(`String`) Name of security intelligence category.
  * `profile_name` - (Optional)(`String`) Name of profile.
* `proxy_profile` - (Optional)(`String`) The proxy profile name.
* `url` - (Optional)(`String`) Configure the url of feed server [https://`<ip or hostname>`:`<port>`/`<uri>`].  
  **Note:** If not set, url is only read from the Junos configuration (so as not to be in conflict with enrollment process).
* `url_parameter` - (Optional)(`String`) Configure the parameter of url.  
  **WARNING** Clear in tfstate.

---

### user_identification arguments

* `ad_access` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to enable 'active-directory-access'. Conflict with `identity_management`.
  * `auth_entry_timeout` - (Optional)(`Int`) Authentication entry timeout number (0, 10-1440) (minutes).
  * `filter_exclude` - (Optional)(`ListOfString`) Exclude addresses.
  * `filter_include` - (Optional)(`ListOfString`) Include addresses.
  * `firewall_auth_forced_timeout` - (Optional)(`Int`) Firewallauth fallback authentication entry forced timeout number (10-1440) (minutes).
  * `invalid_auth_entry_timeout` - (Optional)(`Int`) Invalid authentication entry timeout number (0, 10-1440) (minutes).
  * `no_on_demand_probe` - (Optional)(`bool`) Disable on-demand probe.
  * `wmi_timeout` - (Optional)(`Int`) Wmi timeout number (3..120 seconds).
* `device_info_auth_source` - (Optional)(`String`) Configure authentication-source on device information configuration. Need to be 'active-directory' or 'network-access-controller'.
* `identity_management` - (Optional)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'identity-management' configuration. See the [`identity_management` arguments for user_identification] (#identity_management-arguments-for-user_identification) block. Conflict with `ad_access`.

---

### identity_management arguments for user_identification

* `connection` - (Required)([attribute-as-blocks mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html)) Can be specified only once to declare 'connection' configuration.
  * `primary_address` - (Required)(`String`) IP address of Primary server.
  * `primary_client_id` - (Required)(`String`) Client ID of Primary server for OAuth2 grant.
  * `primary_client_secret` - (Required)(`String`) Client secret of Primary server for OAuth2 grant.  
  **WARNING** Clear in tfstate.
  * `connect_method` - (Optional)(`String`) Method of connection. Need to be 'http' or 'https'.
  * `port` - (Optional)(`Int`) Server port (1..65535).
  * `primary_ca_certificate` - (Optional)(`String`) Ca-certificate file name of Primary server.
  * `query_api` - (Optional)(`String`) Query API.
  * `secondary_address` - (Optional)(`String`) IP address of Secondary server.
  * `secondary_ca_certificate` - (Optional)(`String`) Ca-certificate file name of Secondary server.
  * `secondary_client_id` - (Optional)(`String`) Client ID of Secondary server for OAuth2 grant.
  * `secondary_client_secret` - (Optional)(`String`) Client secret of Secondary server for OAuth2 grant.  
  **WARNING** Clear in tfstate.
  * `token_api` - (Optional)(`String`) API of acquiring token for OAuth2 authentication.
* `authentication_entry_timeout` - (Optional)(`Int`) Authentication entry timeout number (0, 10-1440) (minutes).
* `batch_query_items_per_batch` - (Optional)(`Int`) Items number per batch query (100..1000).
* `batch_query_interval` - (Optional)(`Int`) Query interval for batch query (1..60 seconds).
* `filter_domain` - (Optional)(`ListOfString`) Domain filter.
* `filter_exclude_ip_address_book` - (Optional)(`String`) Referenced address book to exclude IP filter.
* `filter_exclude_ip_address_set` - (Optional)(`String`) Referenced address set to exclude IP filter.
* `filter_include_ip_address_book` - (Optional)(`String`) Referenced address book to include IP filter.
* `filter_include_ip_address_set` - (Optional)(`String`) Referenced address set to include IP filter.
* `invalid_authentication_entry_timeout` - (Optional)(`Int`) Invalid authentication entry timeout number (0, 10-1440) (minutes).
* `ip_query_disable` - (Optional)(`Bool`) Disable IP query.
* `ip_query_delay_time` - (Optional)(`Int`) Delay time to send IP query (0~60sec) (0..60 seconds).
  
## Import

Junos services can be imported using any id, e.g.

```shell
$ terraform import junos_services.services random
```
