resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  passive          = true
  no_install       = true
  no_readvertise   = true
  no_retain        = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop   = "dsc.0"
    preference = 102
    metric     = 102
  }
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination    = "2001:db8:85a3::/48"
  preference     = 100
  metric         = 100
  passive        = true
  no_install     = true
  no_readvertise = true
  no_retain      = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop   = "dsc.0"
    preference = 102
    metric     = 102
  }
  community = ["no-advertise"]
}
