resource "junos_routing_options" "testacc_bgpgroup" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpgroup" {
  name = "testacc_bgpgroup"
  as   = "65000"
}
resource "junos_policyoptions_policy_statement" "testacc_bgpgroup" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_bgpgroup"
  then {
    action = "accept"
  }
}
resource "junos_bgp_group" "testacc_bgpgroup" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup
  ]
  name               = "testacc_bgpgroup"
  routing_instance   = junos_routing_instance.testacc_bgpgroup.name
  advertise_inactive = true
  advertise_peer_as  = true
  as_override        = true
  bgp_multipath {}
  cluster                  = "192.0.2.3"
  damping                  = true
  log_updown               = true
  mtu_discovery            = true
  remove_private           = true
  passive                  = true
  hold_time                = 30
  keep_none                = true
  local_as                 = "65001"
  local_as_private         = true
  local_as_loops           = 1
  local_preference         = 100
  metric_out               = 100
  out_delay                = 30
  peer_as                  = "65002"
  preference               = 100
  authentication_algorithm = "md5"
  local_address            = "192.0.2.3"
  export                   = [junos_policyoptions_policy_statement.testacc_bgpgroup.name]
  import                   = [junos_policyoptions_policy_statement.testacc_bgpgroup.name]
  bfd_liveness_detection {
    detection_time_threshold           = 60
    transmit_interval_threshold        = 30
    transmit_interval_minimum_interval = 10
    holddown_interval                  = 10
    minimum_interval                   = 10
    minimum_receive_interval           = 10
    multiplier                         = 2
    session_mode                       = "automatic"
  }
  family_inet {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  family_inet {
    nlri_type = "multicast"
    accepted_prefix_limit {
      maximum                       = 2
      teardown_idle_timeout_forever = true
    }
    prefix_limit {
      maximum                       = 2
      teardown_idle_timeout_forever = true
    }
  }
  family_inet6 {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  family_inet6 {
    nlri_type = "multicast"
  }
  graceful_restart {
    disable = true
  }
}
