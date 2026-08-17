resource "junos_interface_physical" "testacc_interface_logical_wo_phy" {
  name         = var.interface
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical_wo" {
  name        = "${junos_interface_physical.testacc_interface_logical_wo_phy.name}.100"
  description = "testacc_interface_logical_wo"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
      vrrp_group {
        identifier                    = 100
        virtual_address               = ["192.0.2.2"]
        authentication_type           = "md5"
        authentication_key_wo         = "vrrpKeyOne2"
        authentication_key_wo_version = 2
      }
    }
    address {
      cidr_ip = "192.0.2.129/25"
      vrrp_group {
        identifier                    = 101
        virtual_address               = ["192.0.2.130"]
        authentication_type           = "md5"
        authentication_key_wo         = "vrrpKeyTwo2"
        authentication_key_wo_version = 2
      }
    }
  }
}
