resource "junos_interface_logical" "testacc_v0to1_ipsecvpn" {
  name = "${var.interface}.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_policy" "testacc_v0to1_ipsecvpn" {
  name                = "testacc_v0to1_ipsecvpn"
  proposal_set        = "basic"
  mode                = "main"
  pre_shared_key_text = "thePassWord"
}
resource "junos_security_ike_gateway" "testacc_v0to1_ipsecvpn" {
  name               = "testacc_v0to1_ipsecvpn"
  address            = ["192.0.2.3"]
  policy             = junos_security_ike_policy.testacc_v0to1_ipsecvpn.name
  external_interface = junos_interface_logical.testacc_v0to1_ipsecvpn.name
}
resource "junos_security_ipsec_policy" "testacc_v0to1_ipsecvpn" {
  name         = "testacc_ipsecpol"
  proposal_set = "basic"
  pfs_keys     = "group2"
}
resource "junos_interface_st0_unit" "testacc_v0to1_ipsecvpn" {}
resource "junos_security_ipsec_vpn" "testacc_v0to1_ipsecvpn" {
  name           = "testacc_v0to1_ipsecvpn"
  bind_interface = junos_interface_st0_unit.testacc_v0to1_ipsecvpn.id
  ike {
    gateway          = junos_security_ike_gateway.testacc_v0to1_ipsecvpn.name
    policy           = junos_security_ipsec_policy.testacc_v0to1_ipsecvpn.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  vpn_monitor {
    destination_ip = "192.0.2.129"
    optimized      = true
  }
  establish_tunnels = "on-traffic"
}
