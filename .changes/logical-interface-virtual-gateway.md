<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_interface_logical**:
  * add `virtual_gateway_accept_data`, `virtual_gateway_v4_mac` and `virtual_gateway_v6_mac` arguments
  * add `virtual_gateway_address` argument in `address` block in `family_inet` and `family_inet6` blocks
  * move config errors from defined `vrrp_group` block on `st0.` interface during Plan instead of during Apply
* **data-source/junos_interface_logical**: add `virtual_gateway_accept_data`, `virtual_gateway_v4_mac`, `virtual_gateway_v6_mac` and `virtual_gateway_address` attributes like resource
