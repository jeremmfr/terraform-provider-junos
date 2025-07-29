resource "junos_forwardingoptions_sampling" "testacc_sampInstance" {
  depends_on = [
    junos_interface_logical.testacc_sampInstance,
    junos_system_ntp_server.testacc_sampInstance,
  ]

  routing_instance = junos_routing_instance.testacc_sampInstance.name
  disable          = true
  pre_rewrite_tos  = true
  sample_once      = true
  input {
    rate                   = 1
    max_packets_per_second = 1000
    run_length             = 10
  }
  family_inet_output {
    flow_active_timeout = 60
    flow_server {
      hostname                                              = "192.0.2.1"
      port                                                  = 3000
      version                                               = 8
      aggregation_autonomous_system                         = true
      aggregation_destination_prefix                        = true
      aggregation_protocol_port                             = true
      aggregation_source_destination_prefix                 = true
      aggregation_source_destination_prefix_caida_compliant = true
      aggregation_source_prefix                             = true
      autonomous_system_type                                = "origin"
      source_address                                        = "192.0.2.2"
    }
    interface {
      name           = "si-0/1/0"
      source_address = "192.0.2.2"
      engine_id      = 100
      engine_type    = 100
    }
  }
}
resource "junos_routing_instance" "testacc_sampInstance" {
  name = "testacc_sampInstance"
}

resource "junos_forwardingoptions_sampling" "testacc_sampling" {
  family_inet_input {
    rate                   = 3
    max_packets_per_second = 1000
  }
  family_inet_output {
    flow_inactive_timeout = 60
    flow_server {
      hostname = "192.0.2.11"
      port     = 3002
    }
    flow_server {
      hostname = "192.0.2.10"
      port     = 3003
    }
    interface {
      name           = "si-0/1/0"
      source_address = "192.0.2.2"
    }
  }
}
resource "junos_services_flowmonitoring_v9_template" "testacc_sampling" {
  name = "testacc_sampling"
  type = "ipv4-template"
}

resource "junos_forwardingoptions_sampling" "testacc_sampInstance3" {
  routing_instance = junos_routing_instance.testacc_sampInstance3.name
  family_inet6_input {
    rate                   = 3
    max_packets_per_second = 1000
  }
  family_inet6_output {
    inline_jflow_export_rate    = 4
    inline_jflow_source_address = "192.0.2.2"
    flow_server {
      hostname          = "192.0.2.1"
      port              = 3001
      version9_template = junos_services_flowmonitoring_v9_template.testacc_sampInstance3.name
      local_dump        = true
    }
    flow_server {
      hostname          = "192.0.2.10"
      port              = 3001
      version9_template = junos_services_flowmonitoring_v9_template.testacc_sampInstance3.name
      no_local_dump     = true
    }
  }
}
resource "junos_routing_instance" "testacc_sampInstance3" {
  name = "testacc_sampInstance3"
}
resource "junos_services_flowmonitoring_v9_template" "testacc_sampInstance3" {
  name = "testacc_sampInstance@3"
  type = "ipv6-template"
}

resource "junos_forwardingoptions_sampling" "testacc_sampInstance5" {
  routing_instance = junos_routing_instance.testacc_sampInstance5.name
  family_inet_input {
    rate = 5
  }
  family_inet_output {
    inline_jflow_source_address = "192.0.2.3"
    flow_server {
      hostname          = "192.0.2.10"
      port              = 4000
      version9_template = junos_services_flowmonitoring_v9_template.testacc_sampling.name
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
resource "junos_forwardingoptions_sampling" "testacc_sampInstanceFile" {
  routing_instance = junos_routing_instance.testacc_sampInstanceFile.name
  family_inet_input {
    rate = 2
  }
  family_inet_output {
    file {
      filename       = "flow"
      world_readable = true
      stamp          = true
      size           = 2048
      files          = 5
    }
  }
}
resource "junos_routing_instance" "testacc_sampInstanceFile" {
  name = "testacc_sampInstanceFile"
}
