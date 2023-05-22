<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_bgp_group**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `advertise_external` is now computed to `true` when `advertise_external_conditional` is `true` (instead of the 2 attributes conflicting)
  * `bfd_liveness_detection.version` now generate an error if the value is not in one of strings `0`, `1` or `automatic`
  * add `bgp_error_tolerance`, `description`, `no_client_reflect`, `tcp_aggressive_transmission` arguments
* **resource/junos_bgp_neighbor**:
  * resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
  * `advertise_external` is now computed to `true` when `advertise_external_conditional` is `true` (instead of the 2 attributes conflicting)
  * `bfd_liveness_detection.version` now generate an error if the value is not in one of strings `0`, `1` or `automatic`
  * add `bgp_error_tolerance`, `description`, `no_client_reflect`, `tcp_aggressive_transmission` arguments
