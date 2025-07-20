resource "junos_interface_physical" "testacc_dataIfaceLogInfo" {
  name         = var.interface
  description  = "testacc_dataIfaceLogInfo"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_dataIfaceLogInfo" {
  name        = "${junos_interface_physical.testacc_dataIfaceLogInfo.name}.10"
  description = "testacc_dataIfaceLogInfo"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
    address {
      cidr_ip = "192.0.2.2/25"
    }
  }
  family_inet6 {
    address {
      cidr_ip = "2001:db8::1/64"
    }
  }
}
