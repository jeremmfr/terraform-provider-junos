<!-- markdownlint-disable-file MD013 MD041 -->
# changelog

## v2.0.0 (May 10, 2023)

BREAKING CHANGES with new `v2`:

* **provider**: remove `aes128-cbc` cipher from default ciphers when `ssh_ciphers` is not specified
* remove deprecated **junos_interface** resource and data source
* **resource/junos_bgp_group**, **resource/junos_bgp_neighbor**: remove deprecated `multipath` argument
* **resource/junos_interface_physical**: remove deprecated `ae_lacp`, `ae_link_speed`, `ae_minimum_links` and `ether802_3ad` arguments
* **data-source/junos_interface_physical**: remove same attributes in data source as resource
* **resource/junos_security_ipsec_vpn**: remove deprecated `bind_interface_auto` argument
* **resource/junos_system_radius_server**: remove deprecated `accouting_timeout` attribute

FEATURES:

* add **junos_forwardingoptions_sampling** resource (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))
* add **junos_oam_gretunnel_interface** resource (Fix [#457](https://github.com/jeremmfr/terraform-provider-junos/issues/457))
* add **junos_services_flowmonitoring_v9_template** resource (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))

ENHANCEMENTS:

* refactor provider to integrate new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
the resources and data sources will migrate progressively to this new plugin
* **provider**: add new cipher `aes256-gcm@openssh.com` in default ciphers when `ssh_ciphers` is not specified
* **provider**: enhance debug logs by prefixing messages with local and remote addresses of the sessions
* **resource/junos_forwardingoptions_sampling_instance**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `routing_instance` argument to allow create sampling instance in routing instance.  
  `id` attribute has now the format `<name>_-_<routing_instance>`
  * `flow_server` block argument is now unordered
* **resource/junos_services_flowmonitoring_vipfix_template**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `type` argument now accept `bridge-template`
  * add  `flow_key_output_interface`,  `mpls_template_label_position`, `template_refresh_rate` block argument (Partial fix [#456](https://github.com/jeremmfr/terraform-provider-junos/issues/456))
* **resource/junos_interface_logical**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **data-source/junos_interface_logical**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource  
the data-source schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_interface_physical**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` arguments
* **data-source/junos_interface_physical**:
  * data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource  
  the data-source schema has been upgraded to have one-blocks in single mode instead of list
  * add `encapsulation`, `flexible_vlan_tagging`, `gratuitous_arp_reply`, `hold_time_down`, `hold_time_up`, `link_mode`, `no_gratuitous_arp_reply`, `no_gratuitous_arp_request` and `speed` attributes like resource
* **resource/junos_interface_physical_disable**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_interface_st0_unit**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **data-source/junos_interface_logical_info**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
the data-source schema has been upgraded to have one-blocks in single mode instead of list
* **data-source/junos_interfaces_physical_present**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
* **resource/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value except `type` argument
* **data-source/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource
* **resource/junos_security**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `nat_source` block argument (Fix [#458](https://github.com/jeremmfr/terraform-provider-junos/issues/458))
* **resource/junos_security_address_book**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply (Fix [#403](https://github.com/jeremmfr/terraform-provider-junos/issues/403))  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_global_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_ike_gateway**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `distinguished_name_container` and `distinguished_name_wildcard` arguments inside `remote_identity` block argument
* **resource/junos_security_ike_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` and `reauth_frequency` arguments
* **resource/junos_security_ike_proposal**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_policy**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_proposal**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional string attributes doesn't accept *empty* value
  * add `description` argument
* **resource/junos_security_ipsec_vpn**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * add `copy_outer_dscp`, `manual`, `multi_sa_forwarding_class` and `udp_encapsulate` arguments
* **resource/junos_security_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value  
the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_security_policy_tunnel_pair_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_zone**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
some of config errors are now sent during Plan instead of during Apply  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **data-source/junos_security_zone**: data-source now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource
* **resource/junos_security_zone_book_address**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional boolean attributes doesn't accept value *false*  
optional string attributes doesn't accept *empty* value
* **resource/junos_security_zone_book_address_set**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
optional string attributes doesn't accept *empty* value

BUG FIXES:

* provider: when `ssh_retry_to_establish` > 1, stop retrying to open connections or sessions after a gracefully shutting down call with `Ctrl-c`
* **resource/junos_security_ike_gateway**: don't catch error when read `local_identity` and `remote_identity` block arguments and `type` is `distinguished-name`

## Previous Releases

* [v1](https://github.com/jeremmfr/terraform-provider-junos/blob/v1/CHANGELOG.md)
