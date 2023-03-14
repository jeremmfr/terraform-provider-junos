<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_services_flowmonitoring_vipfix_template**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*)
  * `type` argument now accept `bridge-template`
  * add `flow_key_output_interface` argument
  * add `mpls_template_label_position` argument
  * add `template_refresh_rate` block argument (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))
