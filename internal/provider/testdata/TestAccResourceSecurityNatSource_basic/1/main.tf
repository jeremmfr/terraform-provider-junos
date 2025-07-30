resource "junos_security_nat_source" "testacc_securitySNAT" {
  name        = "testacc_securitySNAT"
  description = "testacc securitySNAT"
  from {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  to {
    type  = "zone"
    value = [junos_security_zone.testacc_securitySNAT.name]
  }
  rule {
    name = "testacc_securitySNATRule"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["tcp"]
    }
    then {
      type = "pool"
      pool = junos_security_nat_source_pool.testacc_securitySNATPool.name
    }
  }
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool" {
  name                                   = "testacc_securitySNATPool"
  description                            = "testacc securitySNATPool"
  address                                = ["192.0.2.1/32", "192.0.2.64/27"]
  routing_instance                       = junos_routing_instance.testacc_securitySNAT.name
  address_pooling                        = "paired"
  port_no_translation                    = true
  pool_utilization_alarm_raise_threshold = 80
  pool_utilization_alarm_clear_threshold = 60
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool2" {
  name             = "testacc_securitySNATPool2"
  description      = "testacc securitySNATPool2"
  address          = ["192.0.2.2/32"]
  routing_instance = junos_routing_instance.testacc_securitySNAT.name
  port_range       = "1300"
}

resource "junos_security_zone" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_routing_instance" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
