---
page_title: "Junos: junos_snmp_v3_usm_user"
---

# junos_snmp_v3_usm_user

Provides a snmp v3 USM user resource.

## Example Usage

```hcl
# Add a snmp v3 usm local-engine user
resource "junos_snmp_v3_usm_user" "user1" {
  name = "user1"
}
# Add a snmp v3 usm remote-engine user
resource "junos_snmp_v3_usm_user" "user2" {
  name        = "user2"
  engine_type = "remote"
  engine_id   = "800007E5804089071BC6D10A41"
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of snmp v3 USM user.
- **engine_type** (Optional, String, Forces new resource)  
  Local or remote engine user.  
  Need to be `local` or `remote`.  
  Defaults to `local`.
- **engine_id** (Optional, String, Forces new resource)  
  Remote engine id (Hex format).
- **authentication_key** (Optional, String, Sensitive)  
  Encrypted key used for user authentication.  
  If the encrypted key is present on Junos device and `authentication_password` is used
  in Terraform config, the value of this argument is left blank to avoid conflict.  
  Conflict with `authentication_password`.
- **authentication_password** (Optional, String, Sensitive)  
  User's authentication password.  
  Due to encryption, when Terraform refreshes the resource, the password can't be read,
  so the provider only checks if it exists and can't detect a change of the password itself
  outside of Terraform.  
  To be able to detect a change of the password outside of Terraform,
  preferably use `authentication_key` argument.  
  Conflict with `authentication_key`.
- **authentication_type** (Optional, String)  
  Define authentication type.  
  Need to be `authentication-md5`, `authentication-none`, `authentication-sha`,
  `authentication-sha224`, `authentication-sha256`, `authentication-sha384` or
  `authentication-sha512`.  
  Defaults to `authentication-none`.  
  `authentication_key` or `authentication_password` need to set when `authentication_type` != `authentication-none`.
- **privacy_key** (Optional, String, Sensitive)  
  Encrypted key used for user privacy.  
  If the encrypted key is present on Junos device and `privacy_password` is used
  in Terraform config, the value of this argument is left blank to avoid conflict.  
  Conflict with `privacy_password`.
- **privacy_password** (Optional, String, Sensitive)  
  User's privacy password.  
  Due to encryption, when Terraform refreshes the resource, the password can't be read,
  so the provider only checks if it exists and can't detect a change of the password itself
  outside of Terraform.  
  To be able to detect a change of the password outside of Terraform,
  preferably use `privacy_key` argument.  
  Conflict with `privacy_key`.
- **privacy_type** (Optional, String)  
  Define privacy type.  
  Need to be `privacy-3des`, `privacy-aes128`, `privacy-des` or `privacy-none`.  
  Defaults to `privacy-none`.  
  `privacy_key` or `privacy_password` need to set when `privacy_type` != `privacy-none`.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `local_-_<name>` or `remote_-_<engine_id>_-_<name>`.

## Import

Junos snmp v3 USM user can be imported using an id made up
of `local_-_<name>` or `remote_-_<engine_id>_-_<name>`, e.g.

```shell
$ terraform import junos_snmp_v3_usm_user.user1 local_-_user1
$ terraform import junos_snmp_v3_usm_user.user2 remote_-_800007E5804089071BC6D10A41_-_user2
```
