resource "junos_routing_instance" "testacc_interface_logical" {
  name = "testacc_interface"
}

resource "junos_interface_logical" "testacc_interface_logical" {
  name = "ip-0/0/0.0"
  tunnel {
    destination                  = "192.0.2.12"
    source                       = "192.0.2.13"
    do_not_fragment              = true
    no_path_mtu_discovery        = true
    routing_instance_destination = junos_routing_instance.testacc_interface_logical.name
    traffic_class                = 202
    ttl                          = 203
  }
}
