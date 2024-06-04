<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_vstp**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*
