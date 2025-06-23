<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_security_utm_custom_url_category**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_custom_url_pattern**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
