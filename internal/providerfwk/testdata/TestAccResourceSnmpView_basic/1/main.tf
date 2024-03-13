
resource "junos_snmp_view" "testacc_snmpview" {
  name        = "testacc_snmpview"
  oid_include = [".1"]
}
