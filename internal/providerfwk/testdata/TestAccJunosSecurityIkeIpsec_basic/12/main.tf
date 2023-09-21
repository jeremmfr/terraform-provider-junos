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
    inet6 = "2001:db8::1"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn2" {
  name            = "testacc_ipsecvpn2"
  copy_outer_dscp = true
  manual {
    external_interface       = junos_interface_logical.testacc_ikegateway.name
    protocol                 = "esp"
    spi                      = 500
    authentication_algorithm = "hmac-sha-256-128"
    authentication_key_hexa  = "00112233445566778899AABBCCDDEEFFaabbccddeeff00112233445566778899"
    encryption_algorithm     = "aes-256-gcm"
    encryption_key_hexa      = "00112233445566778899AABBCCDDEEFFaabbccddeeff00112233445566778899"
    gateway                  = "192.0.2.128"
  }
  multi_sa_forwarding_class = ["network-control", "best-effort"]
  df_bit                    = "clear"
  udp_encapsulate {}
}
