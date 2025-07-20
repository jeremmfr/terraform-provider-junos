resource "junos_access_address_assignment_pool" "testacc_accessAddAssP4" {
  name = "testacc_accessAddAssP4"
  family {
    type    = "inet"
    network = "192.0.2.128/25"
  }
}

resource "junos_routing_instance" "testacc_accessAddAssP6" {
  name = "testacc_accessAddAssP6"
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP6_1" {
  name             = "testacc_accessAddAssP6_1"
  routing_instance = junos_routing_instance.testacc_accessAddAssP6.name
  family {
    type    = "inet6"
    network = "fe80:0:0:b::/64"
  }
}
