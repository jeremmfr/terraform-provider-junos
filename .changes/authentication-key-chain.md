<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add **junos_security_authentication_key_chain** resource

BUG FIXES:

* **resource/junos_security_ike_policy**, **resource/junos_security_ike_proposal**, **resource/junos_security_ipsec_policy**, **resource/junos_security_ipsec_proposal**, **resource/junos_security_ipsec_vpn**: don't check device compatibility with security model (could be used on non-security devices)
