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
resource "junos_interface_logical" "testacc_datainterfaceL2" {
  name                        = "irb.100"
  virtual_gateway_accept_data = true
  virtual_gateway_v4_mac      = "00:aa:bb:cc:dd:ee"
  virtual_gateway_v6_mac      = "00:aa:bb:cc:dd:ff"
  family_inet6 {
    address {
      cidr_ip                 = "fe80::1/64"
      virtual_gateway_address = "fe80::f"
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

data "junos_interface_logical" "testacc_datainterfaceL3" {
  config_interface = "irb.100"
}
