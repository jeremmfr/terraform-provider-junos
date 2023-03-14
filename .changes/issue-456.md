<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_services_flowmonitoring_vipfix_template**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*)
  * `type` argument now accept `bridge-template`
