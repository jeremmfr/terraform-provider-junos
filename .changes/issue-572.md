<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_physical**: add `mc_ae` block argument in `parent_ether_opts` block (Fix [#572](https://github.com/jeremmfr/terraform-provider-junos/issues/572))
* **data-source/junos_interface_physical**: add `mc_ae` block attribute in `parent_ether_opts` block

BUG FIXES:

* **data-source/junos_interface_physical**: fix reading `link_speed` and `minimum_bandwidth` attributes in `parent_ether_opts` block
