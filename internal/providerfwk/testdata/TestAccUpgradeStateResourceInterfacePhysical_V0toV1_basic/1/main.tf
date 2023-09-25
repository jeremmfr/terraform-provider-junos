resource "junos_interface_physical" "testacc_interface" {
  name        = var.interface
  description = "testacc_interfaceU"
  gigether_opts {
    ae_8023ad = var.interfaceAE
  }
}
resource "junos_interface_physical" "testacc_interface2" {
  name        = var.interface2
  description = "testacc_interface2"
  ether_opts {
    flow_control     = true
    loopback         = true
    auto_negotiation = true
  }
  mtu = 9000
}
resource "junos_interface_logical" "testacc_interfaceLO" {
  name = "lo0.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/32"
    }
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [
    junos_interface_physical.testacc_interface,
    junos_interface_logical.testacc_interfaceLO,
  ]
  name        = var.interfaceAE
  description = "testacc_interfaceAE"
  parent_ether_opts {
    bfd_liveness_detection {
      local_address                      = "192.0.2.1"
      detection_time_threshold           = 30
      holddown_interval                  = 30
      minimum_interval                   = 30
      minimum_receive_interval           = 10
      multiplier                         = 1
      neighbor                           = "192.0.2.2"
      no_adaptation                      = true
      transmit_interval_minimum_interval = 10
      transmit_interval_threshold        = 30
      version                            = "automatic"
    }
    no_flow_control   = true
    no_loopback       = true
    link_speed        = "1g"
    minimum_bandwidth = "1 gbps"
  }
  vlan_tagging = true
}
