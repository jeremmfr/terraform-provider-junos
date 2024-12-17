<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_security_idp_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_idp_custom_attack**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `recommended_action` argument is now optional
  * values of `ip_flags` argument in `protocol_ipv4` block has now a config validation to one of `df`, `mf`, `rb`, `no-df`, `no-mf` or `no-rb`
  * values of `tcp_flags` argument in `protocol_tcp` block has now a config validation to one of `ack`, `fin`, `psh`, `r1`, `r2`, `rst`, `syn`, `urg`, `no-acl`, `no-fin`, `no-psh`, `no-r1`, `no-r2`, `no-rst`, `no-syn` or `no-urg`
* **resource/junos_security_idp_custom_attack_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
