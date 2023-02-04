<!-- markdownlint-disable-file MD013 MD041 -->
## upcoming release

BREAKING CHANGES with new `v2`:

* **resource/junos_bgp_group**, **resource/junos_bgp_neighbor**: remove deprecated argument multipath
* **resource/junos_interface_physical**: remove deprecated arguments `ae_lacp`, `ae_link_speed`, `ae_minimum_links` and `ether802_3ad`
* **data-source/junos_interface_physical**: remove same attributes in data source as resource
* **resource/junos_security_ipsec_vpn**: remove deprecated argument `bind_interface_auto`
* **resource/junos_system_radius_server**: remove deprecated attribute `accouting_timeout`
* remove deprecated resource and data source **junos_interface**

FEATURES:

ENHANCEMENTS:

BUG FIXES:

## 1.33.0 (February 07, 2023)

ENHANCEMENTS:

* resource/`junos_interface_physical`: add `mtu` argument (Fixes [#451](https://github.com/jeremmfr/terraform-provider-junos/issues/451))
* data-source/`junos_interface_physical`: add `mtu` attribute (like resource)
* release now with golang 1.20

## 1.32.0 (December 22, 2022)

ENHANCEMENTS:

* provider: add `ssh_retry_to_establish` argument to retry when establishing SSH connections
* resource/`junos_security_address_book`: add `ipv4_only` and `ipv6_only` arguments inside `dns_name` block argument
* refactor the code of most resource reading functions to make it more readable and maintainable
* refactor provider to integrate new plugin [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework),  
  the resources and data sources will migrate progressively to this new plugin

BUG FIXES:

* provider: when `ssh_retry_to_establish` > 1, stop retrying to open connections or sessions after a gracefully shutting down call with `Ctrl-c`

## Previous Releases

* [v1](https://github.com/jeremmfr/terraform-provider-junos/blob/v1/CHANGELOG.md)
