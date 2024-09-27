<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add `junos_applications_ordered` resource, copy of `junos_applications` resource but with Block List instead of Block Set to have a workaround for the performance issue on Block Sets (workaround for [#709](https://github.com/jeremmfr/terraform-provider-junos/issues/709))

ENHANCEMENTS:

BUG FIXES:
