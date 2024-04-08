resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  interface {
    name                           = "all"
    passive                        = true
    authentication_simple_password = "testPass"
    link_protection                = true
    no_advertise_adjacency_segment = true
    no_interface_state_traps       = true
    no_neighbor_down_notification  = true
    poll_interval                  = 19
    te_metric                      = 221
  }
  interface {
    name    = junos_interface_logical.testacc_ospfarea.name
    disable = true
    authentication_md5 {
      key_id = 3
      key    = "testK3y"
    }
    authentication_md5 {
      key_id     = 2
      key        = "testK3y2"
      start_time = "2022-3-9.12:50:00"
    }
    strict_bfd = true
    bfd_liveness_detection {
      minimum_receive_interval           = 29
      transmit_interval_minimum_interval = 48
      transmit_interval_threshold        = 49
      version                            = "automatic"
    }
    neighbor {
      address = "192.0.2.6"
    }
    neighbor {
      address  = "192.0.2.5"
      eligible = "true"
    }
  }
  network_summary_export = [junos_policyoptions_policy_statement.testacc_ospfarea.name]
  network_summary_import = [junos_policyoptions_policy_statement.testacc_ospfarea2.name]
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfareav3ipv4" {
  area_id = "0.0.0.0"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name     = junos_interface_logical.testacc_ospfarea.name
    priority = 0
    bfd_liveness_detection {
      full_neighbors_only                = true
      minimum_receive_interval           = 27
      transmit_interval_minimum_interval = 50
      transmit_interval_threshold        = 51
    }
  }
}
resource "junos_interface_logical" "testacc_ospfarea" {
  name        = "${var.interface}.0"
  description = "testacc_ospfarea"
}
resource "junos_interface_logical" "testacc_ospfarea2" {
  name             = "${var.interface2}.0"
  description      = "testacc_ospfarea2"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
}
resource "junos_routing_instance" "testacc_ospfarea" {
  name = "testacc_ospfarea"
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id          = "0.0.0.0"
  version          = "v3"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
  interface {
    name                = "all"
    passive             = true
    metric              = 100
    retransmit_interval = 32
    hello_interval      = 31
    dead_interval       = 30
    bandwidth_based_metrics {
      bandwidth = "100k"
      metric    = 13
    }
    bandwidth_based_metrics {
      bandwidth = "1m"
      metric    = 14
    }
    demand_circuit                                    = true
    dynamic_neighbors                                 = true
    flood_reduction                                   = true
    interface_type                                    = "p2mp"
    mtu                                               = 900
    no_eligible_backup                                = true
    no_eligible_remote_backup                         = true
    node_link_protection                              = true
    passive_traffic_engineering_remote_node_id        = "192.0.2.7"
    passive_traffic_engineering_remote_node_router_id = "192.0.2.8"
    priority                                          = 21
    transit_delay                                     = 23
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea2.name
    bfd_liveness_detection {
      authentication_loose_check         = true
      detection_time_threshold           = 60
      full_neighbors_only                = true
      holddown_interval                  = 15
      minimum_interval                   = 16
      minimum_receive_interval           = 17
      multiplier                         = 2
      no_adaptation                      = true
      transmit_interval_minimum_interval = 18
      transmit_interval_threshold        = 19
      version                            = "automatic"
    }
  }
}
resource "junos_ospf_area" "testacc_ospfarea2v3realm" {
  area_id          = "0.0.0.0"
  version          = "v3"
  realm            = "ipv4-multicast"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
  interface {
    name    = "all"
    passive = true
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea2.name
    bfd_liveness_detection {
      version                            = "automatic"
      minimum_receive_interval           = 270
      transmit_interval_minimum_interval = 500
      transmit_interval_threshold        = 510
    }
  }
}
