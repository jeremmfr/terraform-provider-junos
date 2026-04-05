resource "junos_forwardingoptions_sampling_instance" "testacc_chassis_fpc" {
  lifecycle {
    create_before_destroy = true
  }

  name = "sampling_instance for testacc_chassis_fpc"
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
  sampling_instance = junos_forwardingoptions_sampling_instance.testacc_chassis_fpc.name
}

resource "junos_chassis_fpc" "testacc_chassis_fpc2" {
  slot_number = 2
  cfp_to_et   = true
}

resource "junos_chassis_fpc" "testacc_chassis_fpc3" {
  slot_number = 3
  error {
    fatal_action    = "log"
    fatal_threshold = 100
    major_action    = "log"
    major_threshold = 70
    minor_action    = "trap"
    minor_threshold = 1
  }
}
