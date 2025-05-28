<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_services_security_intelligence_policy**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional string attributes doesn't accept *empty* value  
