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
    distinguished_name {
      container = "dc=example,dc=com"
    }
    connections_limit = 10
  }
  aaa {
    client_username = "user"
    client_password = "password"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  no_nat_traversal   = true
  dead_peer_detection {
    interval  = 10
    threshold = 3
    send_mode = "probe-idle-tunnel"
  }
  local_address = "192.0.2.4"
}

resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-256-cbc"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name      = "testacc_ipsecpol"
  proposals = [junos_security_ipsec_proposal.testacc_ipsecprop.name]
  pfs_keys  = "group1"
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn2" {
  name = "testacc_ipsecvpn2"
  ike {
    gateway          = junos_security_ike_gateway.testacc_ikegateway.name
    policy           = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  establish_tunnels = "immediately"
  df_bit            = "clear"
}
resource "junos_security_zone" "testacc_secIkeIpsec_local" {
  name = "testacc_secIkeIPsec_local"
  address_book {
    name    = "testacc_vpnlocal"
    network = "192.0.2.64/26"
  }
}
resource "junos_security_zone" "testacc_secIkeIpsec_remote" {
  name = "testacc_secIkeIPsec_remote"
  address_book {
    name    = "testacc_vpnremote"
    network = "192.0.2.128/26"
  }
}
resource "junos_security_policy" "testacc_policyIpsecLocToRem" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_local.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy {
    name                      = "testacc_vpn-out"
    match_source_address      = ["testacc_vpnlocal"]
    match_destination_address = ["testacc_vpnremote"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy" "testacc_policyIpsecRemToLoc" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_remote.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_local.name
  policy {
    name                      = "testacc_vpn-in"
    match_source_address      = ["testacc_vpnremote"]
    match_destination_address = ["testacc_vpnlocal"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy_tunnel_pair_policy" "testacc_vpn-in-out" {
  zone_a        = junos_security_zone.testacc_secIkeIpsec_local.name
  zone_b        = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy_a_to_b = junos_security_policy.testacc_policyIpsecLocToRem.policy[0].name
  policy_b_to_a = junos_security_policy.testacc_policyIpsecRemToLoc.policy[0].name
}
