<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_forwardingoptions_dhcprelay**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list

BUG FIXES:

* **resource/junos_bridge_domain**: fix missing validate of not empty resource in create/update functions
* **resource/junos_ospf_area**: fix missing part of validate config when `version` is null in config (default to `v2`)
