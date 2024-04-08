resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  version = "v3"
  interface {
    name    = "all"
    disable = true
  }
  virtual_link {
    neighbor_id         = "192.0.2.100"
    transit_area        = "192.0.2.101"
    dead_interval       = 102
    demand_circuit      = true
    disable             = true
    flood_reduction     = true
    hello_interval      = 103
    mtu                 = 1040
    retransmit_interval = 105
    transit_delay       = 106

  }
  virtual_link {
    neighbor_id  = "192.0.2.0"
    transit_area = "192.0.2.1"
  }
}
