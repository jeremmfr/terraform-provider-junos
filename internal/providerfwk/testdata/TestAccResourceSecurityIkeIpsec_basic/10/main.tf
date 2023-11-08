resource "junos_interface_logical" "testacc_ikegateway" {
  name = "${var.interface}.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    hostname = "host1.example.com"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-128-cbc"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name      = "testacc_ipsecpol"
  proposals = [junos_security_ipsec_proposal.testacc_ipsecprop.name]
  pfs_keys  = "group2"
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  lifecycle {
    create_before_destroy = true
  }

  name           = "testacc_ipsecvpn"
  bind_interface = junos_interface_logical.testacc_ipsecvpn_bind.name
  ike {
    gateway = junos_security_ike_gateway.testacc_ikegateway.name
    policy  = junos_security_ipsec_policy.testacc_ipsecpol.name
  }
  establish_tunnels = "on-traffic"
  traffic_selector {
    name      = "ts-1"
    local_ip  = "192.0.2.0/26"
    remote_ip = "192.0.3.64/26"
  }
  traffic_selector {
    name      = "ts 2"
    local_ip  = "192.0.2.128/26"
    remote_ip = "192.0.3.192/26"
  }
  udp_encapsulate {
    dest_port = "1025"
  }
}
resource "junos_interface_logical" "testacc_ipsecvpn_bind" {
  name = junos_interface_st0_unit.testacc_ipsec_vpn.id
  family_inet {}
}
resource "junos_interface_st0_unit" "testacc_ipsec_vpn" {}
