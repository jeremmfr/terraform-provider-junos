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

* refactor provider to integrate new plugin [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework),  
  the resources and data sources will migrate progressively to this new plugin
* **resource/junos_security_address_book**: resource now use new plugin [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **resource/junos_security_zone**: resource now use new plugin [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **data-source/junos_security_zone**: resource now use new plugin [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) like resource

BUG FIXES:

* provider: when `ssh_retry_to_establish` > 1, stop retrying to open connections or sessions after a gracefully shutting down call with `Ctrl-c`

## Previous Releases

* [v1](https://github.com/jeremmfr/terraform-provider-junos/blob/v1/CHANGELOG.md)
