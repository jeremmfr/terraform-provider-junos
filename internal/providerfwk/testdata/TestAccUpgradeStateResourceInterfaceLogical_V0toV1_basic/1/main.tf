resource "junos_security_zone" "testacc_interface_logical" {
  name = "testacc_interface"
}
resource "junos_routing_instance" "testacc_interface_logical" {
  name = "testacc_interface"
}
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = var.interface
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  lifecycle {
    create_before_destroy = true
  }
  name                       = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id                    = 101
  description                = "testacc_interface_${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  security_zone              = junos_security_zone.testacc_interface_logical.name
  security_inbound_protocols = ["ospf"]
  security_inbound_services  = ["telnet"]
  routing_instance           = junos_routing_instance.testacc_interface_logical.name
  family_inet {
    mtu = 1500
    rpf_check {
      mode_loose = true
    }
    address {
      cidr_ip   = "192.0.2.1/25"
      primary   = true
      preferred = true
      vrrp_group {
        identifier               = 100
        virtual_address          = ["192.0.2.2"]
        no_accept_data           = true
        advertise_interval       = 10
        advertisements_threshold = 3
        authentication_key       = "thePassWord"
        authentication_type      = "md5"
        no_preempt               = true
        priority                 = 150
      }
    }
  }
  family_inet6 {
    mtu = 1500
    address {
      cidr_ip   = "2001:db8::1/64"
      primary   = true
      preferred = true
      vrrp_group {
        identifier                 = 100
        virtual_address            = ["2001:db8::2"]
        virtual_link_local_address = "fe80::2"
        no_accept_data             = true
        advertise_interval         = 100
        no_preempt                 = true
        priority                   = 150
      }
    }
    address {
      cidr_ip = "fe80::1/64"
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "ip-0/0/0.0"
  tunnel {
    destination                  = "192.0.2.12"
    source                       = "192.0.2.13"
    do_not_fragment              = true
    no_path_mtu_discovery        = true
    routing_instance_destination = junos_routing_instance.testacc_interface_logical.name
    traffic_class                = 202
    ttl                          = 203
  }
}
