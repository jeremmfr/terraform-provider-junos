resource "junos_firewall_filter" "testacc_intlogicalInet" {
  name   = "testacc_intlogicalInet"
  family = "inet"
  term {
    name = "testacc_intlogicalInetTerm"
    then {
      action = "accept"
    }
  }
}
resource "junos_firewall_filter" "testacc_intlogicalInet6" {
  name   = "testacc_intlogicalInet6"
  family = "inet6"
  term {
    name = "testacc_intlogicalInet6Term"
    then {
      action = "accept"
    }
  }
}
resource "junos_security_zone" "testacc_interface_logical" {
  name = "testacc_interface_logical"
}
resource "junos_routing_instance" "testacc_interface_logical" {
  name = "testacc_interface_logical"
}
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = var.interface
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name                       = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  description                = "testacc_interface_${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  disable                    = true
  security_zone              = junos_security_zone.testacc_interface_logical.name
  security_inbound_protocols = ["bgp"]
  security_inbound_services  = ["ssh"]
  routing_instance           = junos_routing_instance.testacc_interface_logical.name
  family_inet {
    mtu           = 1400
    filter_input  = junos_firewall_filter.testacc_intlogicalInet.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet.name
    rpf_check {}
    address {
      cidr_ip = "192.0.2.1/25"
      vrrp_group {
        identifier               = 100
        virtual_address          = ["192.0.2.2"]
        accept_data              = true
        advertise_interval       = 10
        advertisements_threshold = 3
        authentication_key       = "thePassWord"
        authentication_type      = "md5"
        preempt                  = true
        priority                 = 100
        track_interface {
          interface     = junos_interface_physical.testacc_interface_logical_phy.name
          priority_cost = 20
        }
        track_route {
          route            = "192.0.2.128/25"
          routing_instance = "default"
          priority_cost    = 20
        }
      }
    }
  }
  family_inet6 {
    dad_disable   = true
    mtu           = 1400
    filter_input  = junos_firewall_filter.testacc_intlogicalInet6.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet6.name
    address {
      cidr_ip = "2001:db8::1/64"
      vrrp_group {
        identifier                 = 100
        virtual_address            = ["2001:db8::2"]
        virtual_link_local_address = "fe80::2"
        accept_data                = true
        advertise_interval         = 100
        advertisements_threshold   = 3
        preempt                    = true
        priority                   = 100
        track_interface {
          interface     = junos_interface_physical.testacc_interface_logical_phy.name
          priority_cost = 20
        }
        track_route {
          route            = "192.0.2.128/25"
          routing_instance = "default"
          priority_cost    = 20
        }
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
    destination         = "192.0.2.10"
    source              = "192.0.2.11"
    allow_fragmentation = true
    path_mtu_discovery  = true
  }
}
