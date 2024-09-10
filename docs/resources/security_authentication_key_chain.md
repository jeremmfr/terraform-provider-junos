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

## Argument Reference

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of authentication key chain.
- **key** (Required, Block Set)  
  For each identifier `id`, authentication element configuration.
  - **id** (Required, Number)  
    Authentication element identifier.
  - **secret** (Required, String, Sensitive)  
    Authentication key.
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
