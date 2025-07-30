resource "junos_services" "testacc" {
  clean_on_destroy = true
  application_identification {
    no_application_system_cache = true
  }
  user_identification {
    ad_access {
      auth_entry_timeout           = 30
      filter_exclude               = ["192.0.2.3/32", "192.0.2.2/32"]
      filter_include               = ["192.0.2.1/32", "192.0.2.0/32"]
      firewall_auth_forced_timeout = 30
      invalid_auth_entry_timeout   = 30
      no_on_demand_probe           = true
      wmi_timeout                  = 30
    }
  }
}
