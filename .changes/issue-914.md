<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_bgp_group**: add `authentication_key_wo` write-only argument and its `authentication_key_wo_version` companion, to be able to use an ephemeral value for the MD5 authentication key without storing it in the Terraform state (Partial fix [#914](https://github.com/jeremmfr/terraform-provider-junos/issues/914))
* **resource/junos_bgp_neighbor**: add `authentication_key_wo` write-only argument and its `authentication_key_wo_version` companion, to be able to use an ephemeral value for the MD5 authentication key without storing it in the Terraform state (Partial fix [#914](https://github.com/jeremmfr/terraform-provider-junos/issues/914))
* **resource/junos_snmp_v3_usm_user**: add `authentication_key_wo`, `authentication_password_wo`, `privacy_key_wo` and `privacy_password_wo` write-only arguments with their `*_wo_version` companions, to be able to use ephemeral values without storing them in the Terraform state (Partial fix [#914](https://github.com/jeremmfr/terraform-provider-junos/issues/914))
