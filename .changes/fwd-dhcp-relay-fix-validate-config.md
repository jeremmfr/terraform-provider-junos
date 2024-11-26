<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_forwardingoptions_dhcprelay**: fix validation of `attempts` must be in range (1..10) when `version` = v6 in `bulk_leasequery` block
