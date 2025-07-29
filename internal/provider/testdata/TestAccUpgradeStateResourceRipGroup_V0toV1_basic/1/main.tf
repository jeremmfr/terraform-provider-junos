resource "junos_policyoptions_policy_statement" "testacc_ripgroup" {
  lifecycle {
    create_before_destroy = true
  }

  name = "testacc_ripgroup"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ripgroup2" {
  lifecycle {
    create_before_destroy = true
  }

  name = "testacc_ripgroup2"
  then {
    action = "reject"
  }
}
resource "junos_rip_group" "testacc_ripgroup" {
  name           = "test_rip_group#1"
  demand_circuit = true
  bfd_liveness_detection {
    authentication_loose_check         = true
    detection_time_threshold           = 60
    minimum_interval                   = 16
    minimum_receive_interval           = 17
    multiplier                         = 2
    no_adaptation                      = true
    transmit_interval_minimum_interval = 18
    transmit_interval_threshold        = 19
    version                            = "automatic"
  }
  export = [
    junos_policyoptions_policy_statement.testacc_ripgroup.name,
  ]
  import = [
    junos_policyoptions_policy_statement.testacc_ripgroup2.name,
  ]
  max_retrans_time = 101
  metric_out       = 11
  preference       = 1000
  route_timeout    = 90
  update_interval  = 30
}
resource "junos_routing_instance" "testacc_ripgroup2" {
  name = "testacc_ripgroup2"
}
resource "junos_rip_group" "testacc_ripgroup2" {
  name             = "test_rip_group#2"
  routing_instance = junos_routing_instance.testacc_ripgroup2.name
  export = [
    junos_policyoptions_policy_statement.testacc_ripgroup2.name,
    junos_policyoptions_policy_statement.testacc_ripgroup.name,
  ]
  import = [
    junos_policyoptions_policy_statement.testacc_ripgroup2.name,
  ]
}
resource "junos_rip_group" "testacc_ripnggroup" {
  name = "test_ripng_group#1"
  ng   = true

  export = [
    junos_policyoptions_policy_statement.testacc_ripgroup.name,
  ]
  import = [
    junos_policyoptions_policy_statement.testacc_ripgroup2.name,
  ]
  metric_out      = 13
  preference      = 1100
  route_timeout   = 75
  update_interval = 35
}
resource "junos_rip_group" "testacc_ripnggroup2" {
  name             = "test_ripng_group#2"
  ng               = true
  routing_instance = junos_routing_instance.testacc_ripgroup2.name
  export = [
    junos_policyoptions_policy_statement.testacc_ripgroup.name,
  ]
  import = [
    junos_policyoptions_policy_statement.testacc_ripgroup2.name,
    junos_policyoptions_policy_statement.testacc_ripgroup.name,
  ]
}
