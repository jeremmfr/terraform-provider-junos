resource "junos_snmp_clientlist" "testacc_snmpclientlist" {
  name   = "testacc@snmpclientlist"
  prefix = ["192.0.2.1/32", "192.0.2.2/32"]
}
