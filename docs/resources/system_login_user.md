---
page_title: "Junos: junos_system_login_user"
---

# junos_system_login_user

Provides a system login user resource.

## Example Usage

```hcl
# Add a system login user
resource "junos_system_login_user" "user1" {
  name  = "user1"
  class = "operator"
  authentication {
    ssh_public_keys = ["ssh-rsa XXXX user@host"]
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  The name of system login user.
- **class** (Required, String)  
  Login class.
- **uid** (Optional, Computed, Number, Forces new resource)  
  User identifier (uid) (100..64000).
- **authentication** (Optional, Block)  
  Declare `authentication` configuration.
  - **encrypted_password** (Optional, String)  
    Encrypted password string.  
    Conflict with `plain_text_password`.  
    If the encrypted password is present on Junos device and `plain_text_password` is used
    in Terraform config, the value of this argument is left blank to avoid conflict.
  - **no_public_keys** (Optional, Boolean)  
    Disables ssh public key based authentication.
  - **plain_text_password** (Optional, String, Sensitive)  
    Plain text password (auto encrypted by Junos device)  
    Due to encryption, when Terraform refreshes the resource, the plain text password can't be read,
    so the provider only checks if it exists and can't detect a change of the password itself
    outside of Terraform.  
    To be able to detect a change of the password outside of Terraform,
    preferably use `encrypted_password` argument.  
    Conflict with `encrypted_password`.
  - **ssh_public_keys** (Optional, Set of String)  
    Secure shell (ssh) public key string.
- **cli_prompt** (Optional, String)  
  Cli prompt name for this user.
- **full_name** (Optional, String)  
  Full name.

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos system login user can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_system_login_user.user1 user1
```
