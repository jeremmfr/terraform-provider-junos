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
  next_hop = [
    "192.0.2.254"
  ]
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
