<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

ENHANCEMENTS:

* **data-source/junos_interfaces_physical_present**:  
  * add `interfaces` block map attribute with same attributes as `interface_statuses` and additional `logical_interface_names` attribute (Fix [#641](https://github.com/jeremmfr/terraform-provider-junos/issues/641))
  * deprecate `interface_statuses` attribute (read the `interfaces` attribute instead)

BUG FIXES:
