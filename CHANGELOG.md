<!-- markdownlint-disable-file MD013 MD041 -->
## upcoming release

BREAKING CHANGES with new `v2`:

* **resource/junos_bgp_group**, **resource/junos_bgp_neighbor**: remove deprecated `multipath` argument
* **resource/junos_interface_physical**: remove deprecated `ae_lacp`, `ae_link_speed`, `ae_minimum_links` and `ether802_3ad` arguments
* **data-source/junos_interface_physical**: remove same attributes in data source as resource
* **resource/junos_security_ipsec_vpn**: remove deprecated `bind_interface_auto` argument
* **resource/junos_system_radius_server**: remove deprecated `accouting_timeout` attribute
* remove deprecated **junos_interface** resource and data source

FEATURES:

ENHANCEMENTS:

* refactor provider to integrate new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework),  
  the resources and data sources will migrate progressively to this new plugin
* **resource/junos_security_address_book**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **resource/junos_security_ike_gateway**:
  * resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
  * add `distinguished_name_container` and `distinguished_name_wildcard` arguments inside `remote_identity` block argument
* **resource/junos_security_ike_policy**:
  * resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional string attributes doesn't accept *empty* value)
  * add `description` and `reauth_frequency` arguments
* **resource/junos_security_ike_proposal**:
  * resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional string attributes doesn't accept *empty* value)
  * add `description` argument
* **resource/junos_security_policy**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **resource/junos_security_policy_tunnel_pair_policy**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **resource/junos_security_zone**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) and some of config errors are now sent during Plan instead of during Apply (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **data-source/junos_security_zone**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) like resource
* **resource/junos_security_zone_book_address**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value)
* **resource/junos_security_zone_book_address_set**: resource now use new [terraform-plugin-framework](github.com/hashicorp/terraform-plugin-framework) (optional string attributes doesn't accept *empty* value)

BUG FIXES:

* provider: when `ssh_retry_to_establish` > 1, stop retrying to open connections or sessions after a gracefully shutting down call with `Ctrl-c`
* **resource/junos_security_ike_gateway**: don't catch error when read `local_identity` and `remote_identity` block arguments and `type` is `distinguished-name`

## Previous Releases

* [v1](https://github.com/jeremmfr/terraform-provider-junos/blob/v1/CHANGELOG.md)
