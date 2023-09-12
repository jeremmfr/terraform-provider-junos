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
  family_inet6 {
    address {
      cidr_ip = "2001:db8::1/64"
    }
  }
}
