<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

ENHANCEMENTS:

* **resource/junos_vlan**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `community_vlans` argument is now a Set of String (instead of Set of Number) to accept VLAN name in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform
  * `isolated_vlan` argument is now a String (instead of Number) to accept VLAN name in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform
  * `vlan_id` argument is now a String (instead of Number) to accept `all` or `none` in addition to VLAN id  
    data in the state has been updated for the new format  
    Number in config is automatically converted to String by Terraform

BUG FIXES:
