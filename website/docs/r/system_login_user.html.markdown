---
layout: "junos"
page_title: "Junos: junos_system_login_user"
sidebar_current: "docs-junos-resource-system-login-user"
description: |-
  Create a system login user
---

# junos_system_login_user

Provides a system login user resource.

## Example Usage

```hcl
# Add a system login user
resource junos_system_login_user "user1" {
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
  - **no_public_keys** (Optional, Boolean)  
    Disables ssh public key based authentication.
  - **ssh_public_keys** (Optional, Set of String)  
    Secure shell (ssh) public key string.
- **cli_prompt** (Optional, String)  
  Cli prompt name for this user.
- **full_name** (Optional, String)  
  Full name.

## Import

Junos system login user can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_system_login_user.user1 user1
```
