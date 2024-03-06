---
page_title: "Junos: junos_system_root_authentication"
---

# junos_system_root_authentication

-> **Note:** This resource should only be created **once**.
It's used to configure static (not object) options in `system root-authentication` block.  
Destroy this resource has no effect on the Junos configuration.

Configure `system root-authentication` block

## Example Usage

```hcl
# Configure system root-authentication
resource "junos_system_root_authentication" "root_auth" {
  encrypted_password = "$6$XXX"
  ssh_public_keys = [
    "ssh-rsa XXXX",
    "ecdsa-sha2-nistp256 XXXX",
  ]
}
```

## Argument Reference

The following arguments are supported:

-> **Note:** One of `encrypted_password` or `plain_text_password` arguments is required.

- **encrypted_password** (Optional, String)  
  Encrypted password string.  
  If `plain_text_password` is used in Terraform config,
  the value of this argument is left blank to avoid conflict.
- **plain_text_password** (Optional, String, Sensitive)  
  Plain text password (auto encrypted by Junos device)  
  Due to encryption, when Terraform refreshes the resource, the plain text password can't be read,
  so the provider can't detect a change of the password itself outside of Terraform.  
  To be able to detect a change of the password outside of Terraform,
  preferably use `encrypted_password` argument.  
- **no_public_keys** (Optional, Boolean)  
  Disables ssh public key based authentication.
- **ssh_public_keys** (Optional, Set of String)  
  Secure shell (ssh) public key string.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `system_root_authentication`.

## Import

Junos system root-authentication can be imported using any id, e.g.

```shell
$ terraform import junos_system_root_authentication.root_auth random
```
