<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_logical**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value, the resource schema has been upgraded to have one-blocks in single mode instead of list)
* **data-source/junos_interface_logical**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource (schema has been upgraded to have one-blocks in single mode instead of list)
* **resource/junos_interface_physical**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value, the resource schema has been upgraded to have one-blocks in single mode instead of list)
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` arguments
* **data-source/junos_interface_physical**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource (schema has been upgraded to have one-blocks in single mode instead of list)
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` attributes like resource
* **resource/junos_interface_physical_disable**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_interface_st0_unit**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **data-source/junos_interface_logical_info**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (schema has been upgraded to have one-blocks in single mode instead of list)
* **data-source/junos_interfaces_physical_present**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
