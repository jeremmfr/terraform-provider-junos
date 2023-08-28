<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_security_ike_gateway** : fix `Value Conversion Error` when upgrading the provider before `v2.0.0` to `v2.0.0...v2.1.1` and there are this type of resource with `remote_identity` block set in state
