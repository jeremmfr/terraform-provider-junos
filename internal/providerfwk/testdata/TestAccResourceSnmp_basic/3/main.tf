resource "junos_snmp" "testacc_snmp" {
  clean_on_destroy         = true
  arp                      = true
  arp_host_name_resolution = true
  engine_id                = "local \"test#123\""
  health_monitor {}
  routing_instance_access = true
}
