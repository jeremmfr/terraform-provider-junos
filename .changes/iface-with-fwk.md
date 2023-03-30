<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_physical**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value, the resource schema has been upgraded to have one-blocks in single mode instead of list)
* **data-source/junos_interface_physical**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource (schema has been upgraded to have one-blocks in single mode instead of list)
* **resource/junos_interface_physical_disable**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **data-source/junos_interface_logical_info**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (schema has been upgraded to have one-blocks in single mode instead of list)
* **data-source/junos_interfaces_physical_present**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
