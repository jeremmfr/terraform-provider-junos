resource "junos_security" "testacc_security" {
  flow {
    ethernet_switching {
      bypass_non_ip_unicast = true
      no_packet_flooding {
        no_trace_route = true
      }
    }
    tcp_mss {
      all_tcp_mss = 1499
      gre_in {
        mss = 1399
      }
      gre_out {
        mss = 1399
      }
      ipsec_vpn {
        mss = 1399
      }
    }
    tcp_session {
      no_syn_check           = true
      no_syn_check_in_tunnel = true
      time_wait_state {
        apply_to_half_close_state = true
        session_ageout            = true
      }
    }
  }
  idp_sensor_configuration {
    log_suppression {
      include_destination_address = true
    }
    packet_log {
      source_address = "192.0.2.4"
    }
  }
  ike_traceoptions {
    file {
      name              = "ike.log"
      files             = 5
      size              = 100000
      no_world_readable = true
    }
    rate_limit = 100
    # no_remote_trace = true
  }
  log {
    mode                = "event"
    event_rate          = 100
    max_database_record = 1000
    rate_cap            = 100
    source_address      = "192.0.2.1"
  }
  nat_source {
    address_persistent                     = true
    interface_port_overloading_factor      = 32
    pool_default_port_range                = 10242
    pool_default_port_range_to             = 20242
    pool_default_twin_port_range           = 64000
    pool_default_twin_port_range_to        = 65001
    pool_utilization_alarm_clear_threshold = 45
    pool_utilization_alarm_raise_threshold = 80
    port_randomization_disable             = true
    session_drop_hold_down                 = 600
    session_persistence_scan               = true
  }
  policies {
    policy_rematch_extensive = true
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
    feature_profile_web_filtering_juniper_enhanced_server {
      host             = "192.0.2.1"
      port             = 1500
      routing_instance = junos_routing_instance.testacc_security.name
    }
  }
}
resource "junos_routing_instance" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_security"
}
