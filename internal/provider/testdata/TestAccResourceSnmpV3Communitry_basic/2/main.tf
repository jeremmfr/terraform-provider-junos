resource "junos_snmp_v3_community" "testacc_snmpv3comm" {
  community_index = "testacc_snmpv3comm#1"
  security_name   = "testacc_snmpv3comm#1_security2"
  community_name  = "testacc_snmpcomm#1"
  context         = "testacc_snmpv3comm#1_ctx"
  tag             = "testacc_snmpv3comm#1_tag"
}
