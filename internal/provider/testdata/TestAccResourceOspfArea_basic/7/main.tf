resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  version = "v3"
  interface {
    name    = "all"
    disable = true
  }
  virtual_link {
    neighbor_id  = "192.0.2.0"
    transit_area = "192.0.2.1"
  }
}
