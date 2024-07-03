resource "junos_policyoptions_policy_statement" "testacc_ripneigh" {
  lifecycle {
    create_before_destroy = true
  }

  name = "testacc_ripneigh"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ripneigh2" {
  lifecycle {
    create_before_destroy = true
  }

  name = "testacc_ripneigh2"
  then {
    action = "reject"
  }
}
resource "junos_rip_group" "testacc_ripneigh" {
  name = "test_rip_group#1"
}
resource "junos_routing_instance" "testacc_ripneigh2" {
  name = "testacc_ripneigh2"
}
resource "junos_rip_group" "testacc_ripneigh2" {
  name             = "test_rip_group#2"
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_group" "testacc_ripngneigh" {
  name = "test_ripng_group#1"
  ng   = true
}
resource "junos_rip_group" "testacc_ripngneigh2" {
  name             = "test_ripng_group#2"
  ng               = true
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_neighbor" "testacc_ripneigh" {
  name                = "ae0.0"
  group               = junos_rip_group.testacc_ripneigh.name
  any_sender          = true
  authentication_key  = "testKey#1"
  authentication_type = "md5"
  check_zero          = true
  dynamic_peers       = true
  import = [
    junos_policyoptions_policy_statement.testacc_ripneigh.name,
  ]
  interface_type_p2mp = true
  max_retrans_time    = 111
  demand_circuit      = true
  peer = [
    "192.0.2.3",
    "192.0.2.1",
  ]

}
resource "junos_rip_neighbor" "testacc_ripngneigh" {
  name  = "ae0.0"
  ng    = true
  group = junos_rip_group.testacc_ripngneigh.name
  import = [
    junos_policyoptions_policy_statement.testacc_ripneigh2.name,
    junos_policyoptions_policy_statement.testacc_ripneigh.name,
  ]
  metric_in = 3
  receive   = "none"
}
resource "junos_interface_physical" "testacc_ripneigh2" {
  name = var.interface
}
resource "junos_interface_logical" "testacc_ripneigh2" {
  name             = "${junos_interface_physical.testacc_ripneigh2.name}.0"
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
}
resource "junos_rip_neighbor" "testacc_ripneigh2" {
  name             = junos_interface_logical.testacc_ripneigh2.name
  group            = junos_rip_group.testacc_ripneigh2.name
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
  authentication_selective_md5 {
    key_id = 4
    key    = "testKey#4"
  }
  authentication_selective_md5 {
    key_id     = 3
    key        = "testKey#3"
    start_time = "2016-1-1.02:00:00"
  }
  bfd_liveness_detection {
    authentication_loose_check         = true
    detection_time_threshold           = 60
    minimum_interval                   = 16
    minimum_receive_interval           = 17
    multiplier                         = 2
    no_adaptation                      = true
    transmit_interval_minimum_interval = 18
    transmit_interval_threshold        = 19
    version                            = "automatic"
  }
  no_check_zero   = true
  message_size    = 200
  receive         = "both"
  update_interval = 30
  send            = "multicast"
}
resource "junos_rip_neighbor" "testacc_ripngneigh2" {
  name             = junos_interface_logical.testacc_ripneigh2.name
  ng               = true
  group            = junos_rip_group.testacc_ripngneigh2.name
  routing_instance = junos_routing_instance.testacc_ripneigh2.name
  route_timeout    = 300
  send             = "none"
}
