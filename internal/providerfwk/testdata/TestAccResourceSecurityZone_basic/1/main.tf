resource "junos_security_screen" "testaccZone" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "testaccZone"
  description = "testaccZone"
}
resource "junos_security_zone" "testacc_securityZone" {
  name = "testacc_securityZone"
  address_book {
    name        = "testacc_zone1"
    description = "testacc_zone 1"
    network     = "192.0.2.0/25"
  }
  address_book_dns {
    name        = "testacc_zone2"
    description = "testacc_zone 2"
    fqdn        = "test.com"
  }
  address_book_dns {
    name        = "testacc_zone2b"
    description = "testacc_zone 2b"
    fqdn        = "test.com"
    ipv4_only   = true
  }
  address_book_dns {
    name        = "testacc_zone2c"
    description = "testacc_zone 2c"
    fqdn        = "test.com"
    ipv6_only   = true
  }
  address_book_range {
    name        = "testacc_zone3"
    description = "testacc_zone 3"
    from        = "192.0.2.10"
    to          = "192.0.2.12"
  }
  address_book_set {
    name        = "testacc_zoneSet"
    description = "testacc_zone Set"
    address     = ["testacc_zone1"]
  }
  address_book_set {
    name        = "testacc_zoneSet2"
    description = "testacc_zone Set2"
    address     = ["testacc_zone2c"]
    address_set = ["testacc_zoneSet"]
  }
  address_book_wildcard {
    name        = "testacc_zone4"
    description = "testacc_zone 4"
    network     = "192.0.2.0/255.0.255.255"
  }
  application_tracking = true
  inbound_protocols    = ["bgp"]
  description          = "testacc securityZone"
  reverse_reroute      = true
  screen               = junos_security_screen.testaccZone.id
  source_identity_log  = true
  tcp_rst              = true
}
