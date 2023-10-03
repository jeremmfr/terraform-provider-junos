resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = var.interface
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name    = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id = 100
  family_inet {
    dhcp {
      client_identifier_hexadecimal        = "11BBee"
      client_identifier_userid_hexadecimal = "00AAff"
      lease_time_infinite                  = true
    }
  }
  family_inet6 {
    dhcpv6_client {
      client_type                 = "stateful"
      client_identifier_duid_type = "vendor"
      client_ia_type_pd           = true
      req_option                  = ["fqdn", "domain"]
      update_router_advertisement_interface = [
        junos_interface_logical.testacc_interface_logical2.name,
        "${junos_interface_physical.testacc_interface_logical_phy.name}.100",
      ]
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "${junos_interface_physical.testacc_interface_logical_phy.name}.101"
}
