resource "junos_policyoptions_policy_statement" "testacc_ospf" {
  name = "testacc_ospf"
  then {
    action = "accept"
  }
}
resource "junos_ospf" "testacc_ospf" {
  database_protection {
    ignore_count      = 10
    ignore_time       = 600
    maximum_lsa       = 1000
    reset_time        = 700
    warning_only      = true
    warning_threshold = 60
  }
  disable                         = true
  export                          = [junos_policyoptions_policy_statement.testacc_ospf.name]
  external_preference             = 3600
  forwarding_address_to_broadcast = true
  graceful_restart {
    disable             = true
    helper_disable      = true
    helper_disable_type = "both"
    notify_duration     = 900
    restart_duration    = 960
  }
  import               = [junos_policyoptions_policy_statement.testacc_ospf.name]
  labeled_preference   = 5000
  lsa_refresh_interval = 40
  no_nssa_abr          = true
  no_rfc1583           = true
  overload {
    allow_route_leaking = true
    as_external         = true
    stub_network        = true
    timeout             = 600
  }
  preference          = 1000
  prefix_export_limit = 2000
  reference_bandwidth = "1m"
  sham_link           = true
  sham_link_local     = "192.0.2.3"
  spf_options {
    delay                   = 1250
    holddown                = 10500
    no_ignore_our_externals = true
    rapid_runs              = 5
  }
}
resource "junos_routing_instance" "testacc_ospf" {
  name = "testacc_ospf"
}
resource "junos_rib_group" "testacc_ospf_ri" {
  name       = "testacc_ospf_ri"
  import_rib = ["${junos_routing_instance.testacc_ospf.name}.inet.0"]
}
resource "junos_ospf" "testacc_ospf_ri" {
  routing_instance = junos_routing_instance.testacc_ospf.name
  export           = [junos_policyoptions_policy_statement.testacc_ospf.name]
  import           = [junos_policyoptions_policy_statement.testacc_ospf.name]
  domain_id        = "192.0.2.1:100"
  rib_group        = junos_rib_group.testacc_ospf_ri.name
}
resource "junos_ospf" "testacc_ospf_v3" {
  version = "v3"
  export  = [junos_policyoptions_policy_statement.testacc_ospf.name]
  import  = [junos_policyoptions_policy_statement.testacc_ospf.name]
}
