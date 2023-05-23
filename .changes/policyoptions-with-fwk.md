<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_policyoptions_as_path_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  

BUG FIXES:

* **resource/junos_security_ipsec_vpn**: fix length validator (max 31 instead of 32) and remove space exclusion validator of `name` for `traffic_selector` block
