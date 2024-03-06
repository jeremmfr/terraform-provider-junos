---
page_title: "Junos: junos_services"
---

# junos_services

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `services` block.  
By default (without `clean_on_destroy`= true), destroy this resource has no effect on the Junos configuration.

Configure static configuration in `services` block

## Example Usage

```hcl
# Configure services
resource "junos_services" "services" {
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123400"
    url                  = "https://example.com/api/manifest.xml"
  }
}
```

## Argument Reference

The following arguments are supported:

- **clean_on_destroy** (Optional, Boolean)  
  Clean supported lines when destroy this resource.
- **advanced_anti_malware** (Optional, Block)  
  Declare `advanced-anti-malware` static configuration.  
  See [below for nested schema](#advanced_anti_malware-arguments).
- **application_identification** (Optional, Block)  
  Enable `application-identification`.  
  See [below for nested schema](#application_identification-arguments).
- **security_intelligence** (Optional, Block)  
  Declare `security-intelligence` configuration.  
  See [below for nested schema](#security_intelligence-arguments).
- **user_identification** (Optional, Block)  
  Declare `user-identification` configuration.  
  See [below for nested schema](#user_identification-arguments).

---

### advanced_anti_malware arguments

- **connection** (Optional, Block)  
  Declare `connection` configuration.
  - **auth_tls_profile** (Optional, Computed, String)  
    Authentication TLS profile.  
    **Note:** If not set, tls-profile is only read from the Junos configuration
    (so as not to be in conflict with enrollment process).
  - **proxy_profile** (Optional, String)  
    Proxy profile.
  - **source_address** (Optional, String)  
    The source ip for connecting to the cloud server.  
    Conflict with `source_interface`.
  - **source_interface** (Optional, String)  
    The source interface for connecting to the cloud server.  
    Conflict with `source_address`.
  - **url** (Optional, Computed, String)  
    The url of the cloud server [https://`<ip or hostname>`:`<port>`].  
    **Note:** If not set, url is only read from the Junos configuration
    (so as not to be in conflict with enrollment process).
- **default_policy** (Optional, Block)  
  Declare `default-policy` configuration.
  - **blacklist_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware blacklist hit.
  - **default_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware action.
  - **fallback_options_action** (Optional, String)  
    Notification action taken for fallback action.  
    Need to be `block` or `permit`.
  - **fallback_options_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware fallback action.
  - **http_action** (Optional, String)  
    Action taken for contents with verdict meet threshold for HTTP.  
    Need to be `block` or `permit`.  
    `http_inspection_profile` need to be set.
  - **http_client_notify_file** (Optional, String)  
    File name for http response to client notification action taken for contents with verdict meet
    threshold.  
    Conflict with others `http_client_notify_*`.  
    `http_action` and `http_inspection_profile` need to be set.
  - **http_client_notify_message** (Optional, String)  
    Block message to client notification action taken for contents with verdict meet threshold.  
    Conflict with others `http_client_notify_*`.  
    `http_action` and `http_inspection_profile` need to be set.
  - **http_client_notify_redirect_url** (Optional, String)  
    Redirect url to client notification action taken for contents with verdict meet threshold.  
    Conflict with others `http_client_notify_*`.  
    `http_action` and `http_inspection_profile` need to be set.
  - **http_file_verdict_unknown** (Optional, String)  
    Action taken for contents with verdict unknown.  
    `http_action` and `http_inspection_profile` need to be set.
  - **http_inspection_profile** (Optional, String)  
    Advanced Anti-malware inspection-profile name for HTTP.  
    `http_action` need to be set.
  - **http_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware actions for HTTP.  
    `http_action` and `http_inspection_profile` need to be set.
  - **imap_inspection_profile** (Optional, String)  
    Advanced Anti-malware inspection-profile name for IMAP.
  - **imap_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware actions for IMAP.  
    `imap_inspection_profile` need to be set.
  - **smtp_inspection_profile** (Optional, String)  
    Advanced Anti-malware inspection-profile name for SMTP.
  - **smtp_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware actions for SMTP.  
    `smtp_inspection_profile` need to be set.
  - **verdict_threshold** (Optional, String)  
    Verdict threshold.
  - **whitelist_notification_log** (Optional, Boolean)  
    Logging option for Advanced Anti-malware whitelist hit.

---

### application_identification arguments

- **application_system_cache** (Optional, Block)  
  Enable application system cache.  
  Conflict with `no_application_system_cache`.
  - **no_miscellaneous_services** (Optional, Boolean)  
    Disable ASC for miscellaneous services APBR,...
  - **security-services** (Optional, Boolean)  
    Enable ASC for security services (appfw, appqos, idp, skyatp..).
- **no_application_system_cache** (Optional, Boolean)  
  Disable storing AI result in application system cache.  
  Conflict with `application_system_cache`.
- **application_system_cache_timeout** (Optional, Number)  
  Application system cache entry lifetime (0..1000000).
- **download** (Optional, Block)  
  Declare `download` configuration.
  - **automatic_interval** (Optional, Number)  
    Attempt to download new application package (6..720 hours).
  - **automatic_start_time** (Optional, String)  
    Start time to scheduled download and update (MM-DD.hh:mm / YYYY-MM-DD.hh:mm:ss).
  - **ignore_server_validation** (Optional, Boolean)  
    Disable server authentication for Applicaton Signature download.
  - **proxy_profile** (Optional, String)  
    Configure web proxy for Application signature download
  - **url** (Optional, String)  
    URL for application package download.
- **enable_performance_mode** (Optional, Block)  
  Enable performance mode knobs for best DPI performance.
  - **max_packet_threshold** (Optional, Number)  
    Set the maximum packet threshold for DPI performance mode (1..100).
- **global_offload_byte_limit** (Optional, Number)  
  Global byte limit to offload AppID inspection (0..4294967295).
- **imap_cache_size** (Optional, Number)  
  IMAP cache size, it will be effective only after next appid sigpack install (60..512000).
- **imap_cache_timeout** (Optional, Number)  
  IMAP cache entry timeout in seconds (1..86400).
- **inspection_limit_tcp** (Optional, Block)  
  Enable TCP byte/packet inspection limit.
  - **byte_limit** (Optional, Number)  
    TCP byte inspection limit (0..4294967295).
  - **packet_limit** (Optional, Number)  
    TCP packet inspection limit (0..4294967295).
- **inspection_limit_udp** (Optional, Block)  
  Enable UDP byte/packet inspection limit.
  - **byte_limit** (Optional, Number)  
    UDP byte inspection limit (0..4294967295).
  - **packet_limit** (Optional, Number)  
    UDP packet inspection limit (0..4294967295).
- **max_memory** (Optional, Number)  
  Maximum amount of object cache memory JDPI can use (in MB) (1..200000).
- **max_transactions** (Optional, Number)  
  Number of transaction finals to terminate application classification (0..25)
- **micro_apps** (Optional, Boolean)  
  Enable Micro Apps identifcation.
- **statistics_interval** (Optional, Number)  
  Configure application statistics information with collection interval (1..1440 minutes).

---

### security_intelligence arguments

- **authentication_token** (Optional, Computed, String)  
  Token string for authentication to use feed update services.  
  Conflict with `authentication_tls_profile`.  
  **Note:** If not set, token is only read from the Junos configuration
  (so as not to be in conflict with enrollment process).
- **authentication_tls_profile** (Optional, Computed, String)  
  TLS profile for authentication to use feed update services.  
  Conflict with `authentication_token`.  
  **Note:** If not set, tls-profile is only read from the Junos configuration
  (so as not to be in conflict with enrollment process).
- **category_disable** (Optional, Set of String)  
  Categories to be disabled
- **default_policy** (Optional, Block List)  
  For each name of category, configure default-policy for a category.
  - **category_name** (Required, String)  
    Name of security intelligence category.
  - **profile_name** (Required, String)  
    Name of profile.
- **proxy_profile** (Optional, String)  
  The proxy profile name.
- **url** (Optional, Computed, String)  
  Configure the url of feed server [https://`<ip or hostname>`:`<port>`/`<uri>`].  
  **Note:** If not set, url is only read from the Junos configuration
  (so as not to be in conflict with enrollment process).
- **url_parameter** (Optional, String, Sensitive)  
  Configure the parameter of url.  

---

### user_identification arguments

- **ad_access** (Optional, Block)  
  Enable `active-directory-access`.  
  Conflict with `identity_management`.
  - **auth_entry_timeout** (Optional, Number)  
    Authentication entry timeout number (0, 10-1440) (minutes).
  - **filter_exclude** (Optional, Set of String)  
    Exclude addresses.
  - **filter_include** (Optional, Set of String)  
    Include addresses.
  - **firewall_auth_forced_timeout** (Optional, Number)  
    Firewall auth fallback authentication entry forced timeout number (10-1440) (minutes).
  - **invalid_auth_entry_timeout** (Optional, Number)  
    Invalid authentication entry timeout number (0, 10-1440) (minutes).
  - **no_on_demand_probe** (Optional, Boolean)  
    Disable on-demand probe.
  - **wmi_timeout** (Optional, Number)  
    Wmi timeout number (3..120 seconds).
- **device_info_auth_source** (Optional, String)  
  Configure authentication-source on device information configuration.  
  Need to be `active-directory` or `network-access-controller`.
- **identity_management** (Optional, Block)  
  Declare `identity-management` configuration.  
  Conflict with `ad_access`.  
  See [below for nested schema](#identity_management-arguments-for-user_identification).

---

### identity_management arguments for user_identification

- **connection** (Required, Block)  
  Declare `connection` configuration.
  - **primary_address** (Required, String)  
    IP address of Primary server.
  - **primary_client_id** (Required, String)  
    Client ID of Primary server for OAuth2 grant.
  - **primary_client_secret** (Required, String, Sensitive)  
    Client secret of Primary server for OAuth2 grant.  
  - **connect_method** (Optional, String)  
    Method of connection.  
    Need to be `http` or `https`.
  - **port** (Optional, Number)  
    Server port (1..65535).
  - **primary_ca_certificate** (Optional, String)  
    Ca-certificate file name of Primary server.
  - **query_api** (Optional, String)  
    Query API.
  - **secondary_address** (Optional, String)  
    IP address of Secondary server.
  - **secondary_ca_certificate** (Optional, String)  
    Ca-certificate file name of Secondary server.
  - **secondary_client_id** (Optional, String)  
    Client ID of Secondary server for OAuth2 grant.
  - **secondary_client_secret** (Optional, String, Sensitive)  
    Client secret of Secondary server for OAuth2 grant.  
  - **token_api** (Optional, String)  
    API of acquiring token for OAuth2 authentication.
- **authentication_entry_timeout** (Optional, Number)  
  Authentication entry timeout number (0, 10-1440) (minutes).
- **batch_query_items_per_batch** (Optional, Number)  
  Items number per batch query (100..1000).
- **batch_query_interval** (Optional, Number)  
  Query interval for batch query (1..60 seconds).
- **filter_domain** (Optional, Set of String)  
  Domain filter.
- **filter_exclude_ip_address_book** (Optional, String)  
  Referenced address book to exclude IP filter.
- **filter_exclude_ip_address_set** (Optional, String)  
  Referenced address set to exclude IP filter.
- **filter_include_ip_address_book** (Optional, String)  
  Referenced address book to include IP filter.
- **filter_include_ip_address_set** (Optional, String)  
  Referenced address set to include IP filter.
- **invalid_authentication_entry_timeout** (Optional, Number)  
  Invalid authentication entry timeout number (0, 10-1440) (minutes).
- **ip_query_disable** (Optional, Boolean)  
  Disable IP query.
- **ip_query_delay_time** (Optional, Number)  
  Delay time to send IP query (0~60sec) (0..60 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `services`.

## Import

Junos services can be imported using any id, e.g.

```shell
$ terraform import junos_services.services random
```
