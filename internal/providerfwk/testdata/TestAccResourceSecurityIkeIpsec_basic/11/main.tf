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
    inet = "192.168.0.4"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  name = "testacc_ipsecvpn"
  manual {
    external_interface       = junos_interface_logical.testacc_ikegateway.name
    protocol                 = "esp"
    spi                      = 256
    authentication_algorithm = "hmac-sha-256-128"
    authentication_key_text  = "AuthenticationKey123456789012345"
    encryption_algorithm     = "aes-256-gcm"
    encryption_key_text      = "Encryp"
  }
}
