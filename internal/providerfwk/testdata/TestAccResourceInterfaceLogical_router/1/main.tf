resource "junos_interface_logical" "testacc_interface_logical" {
  name = "ip-0/0/0.0"
  tunnel {
    destination         = "192.0.2.10"
    source              = "192.0.2.11"
    allow_fragmentation = true
    path_mtu_discovery  = true
  }
}
resource "junos_interface_logical" "testacc_interface_logical3" {
  name                        = "irb.100"
  virtual_gateway_accept_data = true
  virtual_gateway_v4_mac      = "00:aa:bb:cc:dd:ee"
  virtual_gateway_v6_mac      = "00:aa:bb:cc:dd:ff"
  family_inet {
    address {
      cidr_ip                 = "192.0.2.2/24"
      virtual_gateway_address = "192.0.2.222"
    }

  }
  family_inet6 {
    address {
      cidr_ip                 = "fe80::1/64"
      virtual_gateway_address = "fe80::f"
    }
  }
}
