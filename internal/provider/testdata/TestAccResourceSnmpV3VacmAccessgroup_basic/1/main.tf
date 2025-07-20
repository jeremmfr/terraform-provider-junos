resource "junos_snmp_v3_vacm_accessgroup" "testacc_group" {
  name = "testacc_group#1"
  default_context_prefix {
    model       = "any"
    level       = "none"
    notify_view = "all"
  }
  context_prefix {
    prefix = "ctx#22"
    access_config {
      model     = "any"
      level     = "authentication"
      read_view = "all"
    }
  }
}
resource "junos_snmp_v3_vacm_accessgroup" "testacc_group2" {
  name = "testacc_group#2"
  context_prefix {
    prefix = "ctx"
    access_config {
      model       = "any"
      level       = "none"
      notify_view = "all"
    }
  }
}
