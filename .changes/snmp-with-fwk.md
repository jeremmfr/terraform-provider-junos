<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

ENHANCEMENTS:

* **resource/junos_snmp**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list

BUG FIXES:
