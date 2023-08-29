<!-- markdownlint-disable-file MD013 MD041 -->
BUG FIXES:

* **resource/junos_forwardingoptions_sampling_instance**: avoid resources replacement when upgrading the provider before `v2.0.0` and without refreshing resource states (`-refresh=false`) (Fix [#536](https://github.com/jeremmfr/terraform-provider-junos/issues/536))
