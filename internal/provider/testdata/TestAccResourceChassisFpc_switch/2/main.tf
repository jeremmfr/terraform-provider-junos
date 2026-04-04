resource "junos_forwardingoptions_sampling_instance" "testacc_chassis_fpc_2" {
  name = "sampling_instance for testacc_chassis_fpc #2"
  family_inet_input {
    rate = 2
  }
  family_inet_output {
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname               = "192.0.2.1"
      port                   = 3000
      version_ipfix_template = junos_services_flowmonitoring_vipfix_template.testacc_chassis_fpc.name
    }
  }
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_chassis_fpc" {
  name = "vipfix_template for testacc_chassis_fpc"
  type = "ipv4-template"
}

resource "junos_chassis_fpc" "testacc_chassis_fpc" {
  slot_number       = 0
  sampling_instance = junos_forwardingoptions_sampling_instance.testacc_chassis_fpc_2.name
}

resource "junos_chassis_fpc" "testacc_chassis_fpc2" {
  slot_number = 2
  error {
    fatal_action = "log"
  }
}

resource "junos_chassis_fpc" "testacc_chassis_fpc3" {
  slot_number = 3
  error {
    fatal_action    = "trap"
    fatal_threshold = 99
    major_threshold = 77
    minor_action    = "trap"
    minor_threshold = 32
  }
}
