<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_application**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value
* **data-source/junos_applications**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_application_set**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  * add `application_set`, `description` arguments
* **data-source/junos_application_sets**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
  * add `match_application_sets` argument
  * add `application_set` and `description` attribute in `application_sets` block attribute
