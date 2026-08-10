<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_eventoptions_generate_event**: fix a permanent diff when the `start_time` argument has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
* **resource/junos_ospf_area**: fix a permanent diff when the `start_time` argument in `authentication_md5` block in `interface` block has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
* **resource/junos_rip_neighbor**: fix a permanent diff when the `start_time` argument in `authentication_selective_md5` block has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
* **resource/junos_security**: fix a permanent diff when the `automatic_start_time` argument in `idp_security_package` block has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
* **resource/junos_security_authentication_key_chain**: fix a permanent diff when the `start_time` argument in `key` block has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
* **resource/junos_system_syslog_file**: fix a permanent diff when the `start_time` argument in `archive` block has a leading zero in the month or the day, the device removes it so the value read didn't match the configuration
