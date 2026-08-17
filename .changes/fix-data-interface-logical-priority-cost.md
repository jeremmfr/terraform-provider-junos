<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **data-source/junos_interface_logical**: fix read error when a `vrrp_group` block has a `track_interface` or a `track_route` block, the `priority_cost` attribute was declared as a string instead of a number
