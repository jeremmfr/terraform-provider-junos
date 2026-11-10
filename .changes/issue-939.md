<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_interface_physical**: fix a permanent diff when the `auto_derive_lacp` argument in `esi` block is set, the device returns the statement with its new name `auto-derive lacp-pe-system-id-and-admin-key` (Junos 22.2R1 and later) which was not read (Fix [#939](https://github.com/jeremmfr/terraform-provider-junos/issues/939))
* **data-source/junos_interface_physical**: read `auto_derive_lacp` attribute in `esi` block when the device has the new statement name `auto-derive lacp-pe-system-id-and-admin-key` (Junos 22.2R1 and later)
