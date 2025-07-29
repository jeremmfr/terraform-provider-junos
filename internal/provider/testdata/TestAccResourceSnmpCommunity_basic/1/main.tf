resource "junos_snmp" "testacc_snmpcom" {
  clean_on_destroy        = true
  routing_instance_access = true
}
resource "junos_snmp_community" "testacc_snmpcom" {
  depends_on = [
    junos_snmp.testacc_snmpcom
  ]
  name                    = "testacc_snmpcom@public"
  authorization_read_only = true
  client_list_name        = junos_snmp_clientlist.testacc_snmpcom.name
  routing_instance {
    name = junos_routing_instance.testacc_snmpcom.name
  }
  view = junos_snmp_view.testacc_snmpcom.name
}

resource "junos_snmp_clientlist" "testacc_snmpcom" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_snmpcom"
}
resource "junos_routing_instance" "testacc_snmpcom" {
  name = "testacc_snmpcom"
}
resource "junos_snmp_view" "testacc_snmpcom" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "testacc_snmpcom"
  oid_include = [".1"]
}
