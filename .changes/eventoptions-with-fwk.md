<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

ENHANCEMENTS:

* **resource/junos_eventoptions_destination**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value  
* **resource/junos_eventoptions_generate_event**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
    optional boolean attributes doesn't accept value *false*  
  * add `start_time` argument
* **resource/junos_eventoptions_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list

BUG FIXES: