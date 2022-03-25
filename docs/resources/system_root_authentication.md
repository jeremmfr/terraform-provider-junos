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

- **encrypted_password** (Required, String)  
  Encrypted password string.
- **no_public_keys** (Optional, Boolean)  
  Disables ssh public key based authentication.
- **ssh_public_keys** (Optional, Set of String)  
  Secure shell (ssh) public key string.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with value `system_root_authentication`.

## Import

Junos system root-authentication can be imported using any id, e.g.

```shell
$ terraform import junos_system_root_authentication.root_auth random
```
