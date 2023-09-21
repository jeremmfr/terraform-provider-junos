
resource "junos_security_nat_source" "testacc_securitySNAT" {
  depends_on = [
    junos_security_address_book.testacc_securitySNAT
  ]
  name = "testacc_securitySNAT"
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
      source_address           = ["192.0.2.0/25"]
      source_address_name      = ["testacc_securitySNAT2"]
      source_port              = ["1024", "1021 to 1022"]
      destination_address      = ["192.0.2.128/25"]
      destination_address_name = ["testacc_securitySNAT"]
      destination_port         = ["80", "82 to 83"]
      protocol                 = ["tcp"]
    }
    then {
      type = "pool"
      pool = junos_security_nat_source_pool.testacc_securitySNATPool.name
    }
  }
  rule {
    name = "testacc_securitySNATRule2"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      protocol            = ["udp"]
    }
    then {
      type = "off"
    }
  }
  rule {
    name = "testacc_securitySNATRule3"
    match {
      source_address      = ["192.0.2.0/25"]
      destination_address = ["192.0.2.128/25"]
      application         = ["junos-ssh", "junos-http"]
    }
    then {
      type = "off"
    }
  }
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool" {
  name                    = "testacc_securitySNATPool"
  address                 = ["192.0.2.1/32"]
  routing_instance        = junos_routing_instance.testacc_securitySNAT.name
  address_pooling         = "no-paired"
  port_overloading_factor = 3
}
resource "junos_security_nat_source_pool" "testacc_securitySNATPool2" {
  name             = "testacc_securitySNATPool2"
  description      = "testacc securitySNATPool2"
  address          = ["192.0.2.2/32"]
  routing_instance = junos_routing_instance.testacc_securitySNAT.name
  port_range       = "1300-62000"
}

resource "junos_security_zone" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_routing_instance" "testacc_securitySNAT" {
  name = "testacc_securitySNAT"
}
resource "junos_security_address_book" "testacc_securitySNAT" {
  network_address {
    name  = "testacc_securitySNAT"
    value = "192.0.2.128/27"
  }
  network_address {
    name  = "testacc_securitySNAT2"
    value = "192.0.2.160/27"
  }
}
