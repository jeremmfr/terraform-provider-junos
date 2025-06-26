---
page_title: "Junos: junos_security_utm_policy"
---

# junos_security_utm_policy

Provides a security utm utm-policy resource.

## Example Usage

```hcl
# Add a security utm utm-policy
resource "junos_security_utm_policy" "demo_policy" {
  name = "Demo Policy"
  anti_virus {
    http_profile = "junos-av-defaults"
  }
  traffic_sessions_per_client {
    over_limit = "log-and-permit"
  }
  web_filtering_profile = "junos-wf-local-default"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of security utm utm-policy.
- **anti_spam_smtp_profile** (Optional, String)  
  Name of anti-spam profile.
- **anti_virus** (Optional, Block)  
  Configure for utm anti-virus profile.
  - **ftp_download_profile** (Optional, String)  
    FTP download anti-virus profile.
  - **ftp_upload_profile** (Optional, String)  
    FTP upload anti-virus profile.
  - **http_profile** (Optional, String)  
    HTTP anti-virus profile.
  - **imap_profile** (Optional, String)  
    IMAP anti-virus profile.
  - **pop3_profile** (Optional, String)  
    POP3 anti-virus profile.
  - **smtp_profile** (Optional, String)  
    SMTP anti-virus profile.
- **content_filtering** (Optional, Block)  
  Configure for utm content-filtering profile.
  - **ftp_download_profile** (Optional, String)  
    FTP download content-filtering profile.
  - **ftp_upload_profile** (Optional, String)  
    FTP upload content-filtering profile.
  - **http_profile** (Optional, String)  
    HTTP content-filtering profile.
  - **imap_profile** (Optional, String)  
    IMAP content-filtering profile.
  - **pop3_profile** (Optional, String)  
    POP3 content-filtering profile.
  - **smtp_profile** (Optional, String)  
    SMTP content-filtering profile.
- **content_filtering_rule_set** (Optional, Block List)  
  UTM CF Rule Set.
  - **name** (Required, String)  
    UTM CF Rule-set name.
  - **rule** (Required, Block List)  
    UTM CF Rule.  
    See [below for nested schema](#rule-arguments-for-content_filtering_rule_set).
- **traffic_sessions_per_client** (Optional, Block)  
  Configure for traffic option session per client.
  - **limit** (Optional, Number)  
    Sessions limit.
  - **over_limit** (Optional, String)  
    Over limit action.  
    Need to be `block` or `log-and-permit`.
- **web_filtering_profile** (Optional, String)  
  Web-filtering HTTP profile (local, enhanced, websense).

### rule arguments for content_filtering_rule_set

- **name** (Required, String)  
  UTM CF Rule name.
- **match_applications** (Required, List of String)  
  List of applications to be inspected.  
  Element need to be `any`, `ftp`, `http`, `imap`, `pop3` or `smtp`.
- **match_direction** (Required, String)  
  Direction of the content to be inspected.  
  Need to be `any`, `download` or `upload`.
- **match_file_types** (Required, List of String)  
  List of file-types in match criteria.
- **then_action** (Optional, String)  
  Configure then action.  
  Need to be `block`, `close-client`, `close-client-and-server`, `close-server` or `no-action`.
- **then_notification_log** (Optional, Boolean)  
  Generate security event if content is blocked by rule.
- **then_notification_endpoint** (Optional, Block)  
  Endpoint notification options for the content filtering action taken.
  - **custom_message** (Optional, String)  
    Custom notification message.
  - **notify_mail_sender** (Optional, Boolean)  
    Notify mail sender.  
    `false` to don't notify mail sender.
  - **type** (Optional, String)  
    Endpoint notification type.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm utm-policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_policy.demo_policy "Demo Policy"
```
