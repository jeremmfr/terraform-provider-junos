<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_services**: deprecate `primary_ca_certificate` and `secondary_ca_certificate` arguments in `connection` block in `identity_management` block in `user_identification` block since these options were removed in recent versions of JunOS

* **resource/junos_system**: deprecate `tcp_forwarding` and `no_tcp_forwarding` arguments in `ssh` block in `services` block since these options were removed in recent versions of JunOS

* **resource/junos_system**: add `hostkey_algorithm_list` argument in `ssh` block in `services` block and deprecate `hostkey_algorithm` argument that uses the old JunOS syntax
