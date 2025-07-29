resource "junos_forwardingoptions_sampling_instance" "testacc_sampInstance" {
  depends_on = [
    junos_interface_logical.testacc_sampInstance,
    junos_system_ntp_server.testacc_sampInstance,
  ]

  name = "testacc_instance@1"
  input {
    rate = 1
  }
  family_inet_output {
    flow_server {
      hostname = "192.0.2.1"
      port     = 3000
    }
    interface {
      name           = "si-0/1/0"
      source_address = "192.0.2.2"
    }
  }
}
resource "junos_forwardingoptions_sampling_instance" "testacc_sampInstance2" {
  name = "testacc_instance@2"
  family_inet_input {
    rate = 2
  }
  family_inet_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = junos_services_flowmonitoring_vipfix_template.testacc_sampInstance2.name
    }
  }
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_sampInstance2" {
  name = "testacc_sampInstance@2"
  type = "ipv4-template"
}

resource "junos_forwardingoptions_sampling_instance" "testacc_sampInstance3" {
  name = "testacc_instance@3"
  family_inet6_input {
    rate = 2
  }
  family_inet6_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = junos_services_flowmonitoring_vipfix_template.testacc_sampInstance3.name
    }
  }
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_sampInstance3" {
  name = "testacc_sampInstance@3"
  type = "ipv6-template"
}

resource "junos_forwardingoptions_sampling_instance" "testacc_sampInstance4" {
  name = "testacc_instance@4"
  family_mpls_input {
    rate = 2
  }
  family_mpls_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = junos_services_flowmonitoring_vipfix_template.testacc_sampInstance4.name
    }
  }
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_sampInstance4" {
  name = "testacc_sampInstance@4"
  type = "mpls-template"
}

resource "junos_forwardingoptions_sampling_instance" "testacc_sampInstance5" {
  name             = "testacc_instance@5"
  routing_instance = junos_routing_instance.testacc_sampInstance5.name
  family_inet_input {
    rate = 2
  }
  family_inet_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = junos_services_flowmonitoring_vipfix_template.testacc_sampInstance2.name
    }
  }
}
resource "junos_routing_instance" "testacc_sampInstance5" {
  name = "testacc_sampInstance5"
}

resource "junos_system_ntp_server" "testacc_sampInstance" {
  address = "192.0.2.3"
}
resource "junos_interface_logical" "testacc_sampInstance" {
  name = "si-0/1/0.0"
  family_inet {}
}
