resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1", ".1.1"]
  oid_exclude = [".1.1.2"]
}
