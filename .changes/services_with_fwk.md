<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_services**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)  
  some of config errors are now sent during Plan instead of during Apply  
  optional boolean attributes doesn't accept value *false*  
  optional string attributes doesn't accept *empty* value  
  the resource schema has been upgraded to have one-blocks in single mode instead of list  
  computed attributes, in `advanced_anti_malware.connection` and `security_intelligence` block, are now unknown if block is specified without these attributes when updating resources (except if attributes are null in Terraform state), instead of using values in Terraform state
