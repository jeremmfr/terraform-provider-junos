<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_group_dual_system**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false* (but still accepted for `apply_groups` argument)  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
