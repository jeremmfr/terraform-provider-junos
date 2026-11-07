resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  active           = true
  install          = true
  readvertise      = true
  no_resolve       = true
  retain           = true
  bfd_liveness_detection {
    authentication_loose_check         = true
    detection_time_threshold           = 60
    holddown_interval                  = 1000
    minimum_interval                   = 16
    minimum_receive_interval           = 17
    multiplier                         = 2
    neighbor                           = "192.0.2.253"
    no_adaptation                      = true
    transmit_interval_minimum_interval = 18
    transmit_interval_threshold        = 19
    version                            = "automatic"
  }
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop  = "192.0.2.250"
    interface = "st0.0"
  }
  community = ["no-advertise"]
}
resource "junos_static_route" "testacc_staticRoute_default" {
  destination = "192.0.2.0/25"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  active      = true
  install     = true
  readvertise = true
  no_resolve  = true
  retain      = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  community                    = ["no-advertise"]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_static_route" "testacc_staticRoute2_default" {
  destination = "192.0.2.128/28"
  resolve     = true # without it, vSRX refuses local_address in bfd_liveness_detection (single-hop session)
  next_hop = [
    "192.0.2.254"
  ]
  bfd_liveness_detection {
    local_address       = "192.0.2.1"
    minimum_interval    = 1000
    minimum_receive_ttl = 10
  }
}
resource "junos_static_route" "testacc_staticRoute3_default" {
  destination = "192.0.2.144/28"
  resolve     = true # without it, vSRX refuses local_address in bfd_liveness_detection (single-hop session)
  qualified_next_hop {
    next_hop = "192.0.2.254"
    bfd_liveness_detection {
      local_address       = "192.0.2.1"
      minimum_interval    = 1000
      minimum_receive_ttl = 10
    }
  }
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination = "2001:db8:85a3::/48"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  active      = true
  install     = true
  readvertise = true
  no_resolve  = true
  retain      = true
  bfd_liveness_detection {
    minimum_interval = 1000
    neighbor         = "2001:db8:85a4::2"
  }
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  community                    = ["no-advertise"]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance" {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  active           = true
  install          = true
  readvertise      = true
  no_resolve       = true
  retain           = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop  = "2001:db8:85a4::1"
    interface = "st0.0"
  }
  community = ["no-advertise"]
}
