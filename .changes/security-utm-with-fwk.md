<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_security_utm_custom_url_category**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_custom_url_pattern**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
* **resource/junos_security_utm_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `content_filtering_rule_set` block argument
* **resource/junos_security_utm_profile_web_filtering_juniper_enhanced**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_utm_profile_web_filtering_juniper_local**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `no_safe_search` argument
