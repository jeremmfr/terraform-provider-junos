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
- **traffic_sessions_per_client** (Optional, Block)  
  Configure for traffic option session per client.
  - **limit** (Optional, Number)  
    Sessions limit.
  - **over_limit** (Optional, String)  
    Over limit action
- **web_filtering_profile** (Optional, String)  
  Web-filtering HTTP profile (local, enhanced, websense)

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security utm utm-policy can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_utm_policy.demo_policy "Demo Policy"
```
