<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_access_address_assignment_pool**: 
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `user_name` argument inside `host` block inside `family` block (`hardware_address` is now optional)
  * `host` block, `xauth_attributes_primary_dns` and `xauth_attributes_secondary_dns` arguments inside `family` block can now be configured when `type` = `inet6`
