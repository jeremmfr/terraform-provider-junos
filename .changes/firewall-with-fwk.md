<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_firewall_policer**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `logical_bandwidth_policer`, `logical_interface_policer`, `physical_interface_policer`, `shared_bandwidth_policer` and `if_exceeding_pps` arguments
