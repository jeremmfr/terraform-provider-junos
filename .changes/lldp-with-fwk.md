<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_lldp_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list  
* **resource/junos_lldpmed_interface**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list  
