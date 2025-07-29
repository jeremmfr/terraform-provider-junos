resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  interface {
    name                = "all"
    disable             = true
    passive             = true
    metric              = 100
    retransmit_interval = 12
    hello_interval      = 11
    dead_interval       = 10
  }
  interface {
    name      = junos_interface_logical.testacc_ospfarea.name
    secondary = true
  }
}
resource "junos_ospf_area" "testacc_ospfareav3ipv4" {
  area_id = "0"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name    = "all"
    disable = true
  }
  interface {
    name      = junos_interface_logical.testacc_ospfarea.name
    secondary = true
  }
}
resource "junos_interface_logical" "testacc_ospfarea" {
  name        = "${var.interface}.0"
  description = "testacc_ospfarea"
}

