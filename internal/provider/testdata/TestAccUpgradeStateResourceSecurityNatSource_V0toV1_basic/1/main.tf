resource "junos_security_nat_source" "testacc_securitySNAT" {
  name        = "testacc_securitySNAT_upgrade"
  description = "testacc securitySNAT upgrade"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT_upgrade.name]
  }
  to {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT_upgrade.name]
  }
  rule {
    name = "testacc_securitySNATRule"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["tcp"]
    }
    then {
      type = "off"
    }
  }
}
resource "junos_security_zone" "testacc_securitySNAT_upgrade" {
  name = "testacc_securitySNAT_upgrade"
}
