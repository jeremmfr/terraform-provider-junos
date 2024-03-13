resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp" {
  name  = "testacc_snmpv3secutogrp"
  model = "usm"
  group = "testacc_snmpv3secutogrp2"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp2" {
  name  = "testacc_snmpv3secutogrp"
  model = "v1"
  group = "testacc_snmpv3secutogrp2"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp3" {
  name  = "testacc_snmpv3secutogrp"
  model = "v2c"
  group = "testacc_snmpv3secutogrp2"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp4" {
  name  = " testacc_snmpv3secutogrp#4"
  model = "usm"
  group = " testacc_snmpv3secutogrp#2"
}
