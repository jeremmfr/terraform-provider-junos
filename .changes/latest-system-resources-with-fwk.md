<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_system_ntp_server**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attribute doesn't accept value *false*  
  optional string attribute doesn't accept *empty* value  
  * add `nts` block argument
* **resource/junos_system_services_dhcp_localserver_group**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list

BUG FIXES:

* **resource/junos_forwardingoptions_dhcprelay**:
  * fix missing detection of conflict between `dynamic_profile_aggregate_clients` and `dynamic_profile_use_primary` arguments in config validation
  * fix missing detection of empty `overrides_v4` and `overrides_v6` block arguments in config validation
* **resource/junos_forwardingoptions_dhcprelay_group**:
  * fix missing detection of conflict between `dynamic_profile_aggregate_clients` and `dynamic_profile_use_primary` arguments in config validation
  * fix missing detection of empty `overrides_v4` and `overrides_v6` block arguments in config validation
