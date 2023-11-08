resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = var.interface
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_datainterfaceL" {
  name        = "${junos_interface_physical.testacc_datainterfaceP.name}.100"
  description = "testacc_datainterfaceL"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
}

data "junos_interface_logical" "testacc_datainterfaceL" {
  config_interface = var.interface
  match            = "192.0.2.1/"
}

data "junos_interface_logical" "testacc_datainterfaceL2" {
  match = "192.0.2.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
}
