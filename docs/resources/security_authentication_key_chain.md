---
page_title: "Junos: junos_security_authentication_key_chain"
---

# junos_security_authentication_key_chain

Provides a security authentication key chain resource.

## Example Usage

```hcl
# Add an authentication key chain
resource "junos_security_authentication_key_chain" "demo" {
  name = "chain1"
  key {
    id         = 0
    secret     = "aS3cret#1"
    start_time = "2021-12-11.10:09:08"
  }
}
```

```hcl
# Add an authentication key chain with the keys not stored in the Terraform state
resource "junos_security_authentication_key_chain" "demo_wo" {
  name = "chain2"
  key {
    id         = 0
    start_time = "2021-12-11.10:09:08"
  }
  key {
    id         = 1
    start_time = "2022-06-01.00:00:00"
  }
  key_secret_wo = {
    "0" = {
      value   = ephemeral.example.demo.secret0
      version = 1
    }
    "1" = {
      value   = ephemeral.example.demo.secret1
      version = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

-> **Note**
  For each `key` block, one of the `secret` argument or an entry with the same `id` in the
  `key_secret_wo` argument is required.

- **name** (Required, String, Forces new resource)  
  Name of authentication key chain.
- **key** (Required, Block Set)  
  For each identifier `id`, authentication element configuration.
  - **id** (Required, Number)  
    Authentication element identifier.
  - **secret** (Optional, String, Sensitive)  
    Authentication key.  
    Conflict with an entry with the same `id` in `key_secret_wo`.
  - **start_time** (Required, String)  
    Start time for key transmission (YYYY-MM-DD.HH:MM:SS).
  - **algorithm** (Optional, String)  
    Authentication algorithm.
  - **ao_cryptographic_algorithm** (Optional, String)  
    Cryptographic algorithm for TCP-AO Traffic key and MAC digest generation.
  - **ao_recv_id** (Optional, Number)  
    Recv id for TCP-AO entry (0..255).
  - **ao_send_id** (Optional, Number)  
    Send id for TCP-AO entry (0..255).
  - **ao_tcp_ao_option** (Optional, String)  
    Include TCP-AO option within message header.  
    Need to be `disabled` or `enabled`.
  - **key_name** (Optional, String)  
    Key name in hexadecimal format used for macsec.
  - **options** (Optional, Sstring)  
    Protocol's transmission encoding format.  
    Need to be `basic` or `isis-enhanced`.
- **description** (Optional, String)  
  Text description of this authentication-key-chain.
- **key_secret_wo** (Optional, Map of Block)  
  For each authentication element identifier, authentication key not stored in state.  
  The key of each entry is the `id` of the `key` block to which it applies.  
  Requires Terraform 1.11 or later.
  - **value** (Required, String, Sensitive, Write-only)  
    Authentication key.  
    Conflict with the `secret` argument of the `key` block with the same `id`.
  - **version** (Required, Number)  
    Version of `value` to trigger the sending of its value.  
    Increment it to send the current value of `value` to the device.
- **tolerance** (Optional, Number)  
  Clock skew tolerance (0..4294967295 seconds).

## Attribute Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos security authentication key chain can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_security_authentication_key_chain.demo chain1
```

!> **Warning**
  Write-only arguments cannot be filled by an import, so the authentication keys read on the device
  are stored in the `secret` argument of each `key` block, and therefore in the Terraform state.  
  When the configuration uses `key_secret_wo`, the next apply removes them from the state.
