<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_firewall_filter**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `destination_mac_address`, `destination_mac_address_except`,
  `forwarding_class`, `forwarding_class_except`,
  `interface`,
  `loss_priority`, `loss_priority_except`,
  `packet_length`, `packet_length_except`,
  `policy_map`, `policy_map_except`,
  `source_mac_address` and `source_mac_address_except` arguments in from block in term block
  * add `forwarding_class` and `loss_priority` arguments in then block in term block
* **resource/junos_firewall_policer**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    some of config errors are now sent during Plan instead of during Apply  
    optional boolean attributes doesn't accept value *false*  
    optional string attributes doesn't accept *empty* value  
    the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `logical_bandwidth_policer`, `logical_interface_policer`, `physical_interface_policer`, `shared_bandwidth_policer` and `if_exceeding_pps` arguments
