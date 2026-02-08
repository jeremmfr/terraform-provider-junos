<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_logical**: add `filter_input_list` and `filter_output_list` arguments in `family_inet` and `family_inet6` block (Fix [#865](https://github.com/jeremmfr/terraform-provider-junos/issues/865))
* **data-source/junos_interface_logical**: add `filter_input_list` and `filter_output_list` attributes in `family_inet` and `family_inet6` block like resource
