data "junos_system_information" "srx" {}
locals {
  ipsec_vpn_manual_available = tonumber(replace(data.junos_system_information.srx.os_version, "/\\..*$/", "")) < 22 ? 1 : 0
}

resource "junos_interface_logical" "testacc_ipsecvpn_wo" {
  name = "${var.interface}.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn_wo" {
  count = local.ipsec_vpn_manual_available

  name = "testacc_ipsecvpn_wo"
  manual {
    external_interface                 = junos_interface_logical.testacc_ipsecvpn_wo.name
    protocol                           = "esp"
    spi                                = 256
    authentication_algorithm           = "hmac-sha-256-128"
    authentication_key_hexa_wo         = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
    authentication_key_hexa_wo_version = 1
    encryption_algorithm               = "aes-256-gcm"
    encryption_key_text                = "Encryp"
  }
}
