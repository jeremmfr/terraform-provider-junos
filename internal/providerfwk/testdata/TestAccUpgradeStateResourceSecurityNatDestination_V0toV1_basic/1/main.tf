resource "junos_security_nat_destination" "testacc_securityDNAT" {
  name        = "testacc_securityDNAT_upgrade"
  description = "testacc securityDNAT upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securityDNAT_upgrade.name]
  }
  rule {
    name                = "testacc_securityDNATRule"
    destination_address = "192.0.2.1/32"
    then {
      type = "off"
    }
  }
}
resource "junos_security_zone" "testacc_securityDNAT_upgrade" {
  name = "testacc_securityDNAT_upgrade"
}
