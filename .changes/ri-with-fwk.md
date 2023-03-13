<!-- markdownlint-disable-file MD013 MD041 -->
ENHANCEMENTS:

* **resource/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (optional boolean attributes doesn't accept value *false*, optional string attributes doesn't accept *empty* value except `type` argument)
* **data-source/junos_routing_instance**: resource now use new [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) like resource
