resource "junos_snmp" "testacc_snmpcom" {
  clean_on_destroy        = true
  routing_instance_access = true
}
resource "junos_snmp_community" "testacc_snmpcom" {
  depends_on = [
    junos_snmp.testacc_snmpcom
  ]
  name                     = "testacc_snmpcom@public"
  authorization_read_write = true
  clients                  = ["192.0.2.0/24"]
  routing_instance {
    name = junos_routing_instance.testacc_snmpcom.name
  }
  routing_instance {
    name = junos_routing_instance.testacc2_snmpcom.name
  }
}

resource "junos_routing_instance" "testacc_snmpcom" {
  name = "testacc_snmpcom"
}
resource "junos_routing_instance" "testacc2_snmpcom" {
  name = "testacc2_snmpcom"
}
