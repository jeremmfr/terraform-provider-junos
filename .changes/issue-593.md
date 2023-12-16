<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add `junos_system_syslog_user` resource (Fix [#593](https://github.com/jeremmfr/terraform-provider-junos/issues/593))

ENHANCEMENTS:

* **resource/junos_system_syslog_file**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list
* **resource/junos_system_syslog_host**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list

BUG FIXES:

* **resource/junos_system_syslog_file**: fix reading `archive size` when value is a multiple of 1024 (k,m,g)
