resource "junos_policyoptions_policy_statement" "testacc_ospf" {
  name = "testacc_ospf"
  then {
    action = "accept"
  }
}
resource "junos_ospf" "testacc_ospf" {
  database_protection {
    maximum_lsa = 10
  }
  disable                         = true
  export                          = [junos_policyoptions_policy_statement.testacc_ospf.name]
  external_preference             = 3600
  forwarding_address_to_broadcast = true
  graceful_restart {
    disable                = true
    no_strict_lsa_checking = true
    notify_duration        = 900
    restart_duration       = 960
  }
  import               = [junos_policyoptions_policy_statement.testacc_ospf.name]
  labeled_preference   = 5000
  lsa_refresh_interval = 40
  no_nssa_abr          = true
  no_rfc1583           = true
  overload {}
  preference          = 1000
  prefix_export_limit = 2000
  reference_bandwidth = "10k"
  sham_link           = true
  sham_link_local     = "192.0.2.3"
  spf_options {
    delay                   = 1250
    holddown                = 10500
    no_ignore_our_externals = true
    rapid_runs              = 5
  }
}
