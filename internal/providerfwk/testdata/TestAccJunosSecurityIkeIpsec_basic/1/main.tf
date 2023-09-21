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
  dh_group                 = "group2"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "main"
  pre_shared_key_text = "thePassWord"
  reauth_frequency    = 50
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name               = "testacc_ikegateway"
  address            = ["192.0.2.3"]
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  general_ike_id     = true
  no_nat_traversal   = true
  dead_peer_detection {
    interval  = 10
    threshold = 3
    send_mode = "always-send"
  }
  local_address = "192.0.2.4"
  local_identity {
    type  = "hostname"
    value = "testacc"
  }
  version = "v2-only"
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
resource "junos_interface_st0_unit" "testacc_ipsecvpn" {}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  name           = "testacc_ipsecvpn"
  bind_interface = junos_interface_st0_unit.testacc_ipsecvpn.id
  ike {
    gateway          = junos_security_ike_gateway.testacc_ikegateway.name
    policy           = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  vpn_monitor {
    destination_ip        = "192.0.2.129"
    optimized             = true
    source_interface_auto = true
  }
  establish_tunnels = "on-traffic"
  df_bit            = "clear"
}
