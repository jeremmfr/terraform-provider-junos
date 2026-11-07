<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_logical**: add `tcp_mss` argument in `family_inet` and `family_inet6` block (Fix [#935](https://github.com/jeremmfr/terraform-provider-junos/issues/935))
* **data-source/junos_interface_logical**: add `tcp_mss` attribute in `family_inet` and `family_inet6` block like resource
