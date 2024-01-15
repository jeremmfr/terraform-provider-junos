<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_aggregate_route, junos_application, junos_bgp_group, junos_bgp_neighbor, junos_bridge_domain, junos_interface_physical**:  
avoid trigger the conflict errors when Terraform call resource config validate and value for potential conflict is unknown (can be null afterwards) (Fix [#611](https://github.com/jeremmfr/terraform-provider-junos/issues/611))
