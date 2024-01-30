resource "junos_snmp" "testacc_snmp" {
  arp                        = true
  contact                    = "contact@example.com"
  description                = "snmp description"
  engine_id                  = "use-mac-address"
  filter_duplicates          = true
  filter_interfaces          = ["(ge|xe|ae).*\\.0", "fxp0"]
  filter_internal_interfaces = true
  health_monitor {
    falling_threshold     = 41
    idp                   = true
    idp_falling_threshold = 42
    idp_interval          = 43
    idp_rising_threshold  = 44
    interval              = 45
    rising_threshold      = 46
  }
  if_count_with_filter_interfaces = true
  interface                       = ["fxp0.0"]
  location                        = "Paris, France"
  routing_instance_access         = true
  routing_instance_access_list    = [junos_routing_instance.testacc_snmp.name]
}

resource "junos_routing_instance" "testacc_snmp" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_snmp"
}
