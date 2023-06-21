<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_application_set**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  * add `application_set`, `description` arguments
* **data-source/junos_application_sets**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
  * add `match_application_sets` argument
  * add `application_set` and `description` attribute in `application_sets` block attribute
