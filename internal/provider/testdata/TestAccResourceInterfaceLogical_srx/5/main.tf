resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = var.interface
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name    = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id = 100
  family_inet {
    dhcp {
      client_identifier_ascii                        = "BBAA#1"
      client_identifier_prefix_hostname              = true
      client_identifier_prefix_routing_instance_name = true
      client_identifier_use_interface_description    = "device"
      client_identifier_userid_ascii                 = "BBCC#2"
      force_discover                                 = true
      lease_time                                     = 600
      metric                                         = 0
      no_dns_install                                 = true
      options_no_hostname                            = true
      retransmission_attempt                         = 0
      retransmission_interval                        = 4
      server_address                                 = "192.0.2.1"
      update_server                                  = true
      vendor_id                                      = 2
    }
  }
  family_inet6 {
    dhcpv6_client {
      client_type                               = "stateful"
      client_identifier_duid_type               = "vendor"
      client_ia_type_na                         = true
      client_ia_type_pd                         = true
      no_dns_install                            = true
      prefix_delegating_preferred_prefix_length = 0
      prefix_delegating_sub_prefix_length       = 5
      rapid_commit                              = true
      req_option                                = ["fqdn"]
      retransmission_attempt                    = 0
      update_router_advertisement_interface = [
        junos_interface_logical.testacc_interface_logical2.name,
      ]
      update_server = true
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name          = "${junos_interface_physical.testacc_interface_logical_phy.name}.101"
  encapsulation = "dix"
}
