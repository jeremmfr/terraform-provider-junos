resource "junos_interface_logical" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "${var.interface}.0"
  description = "testacc_security"
}
resource "junos_services_proxy_profile" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name               = "testacc_security"
  protocol_http_host = "192.0.2.11"
  protocol_http_port = 3128
}
resource "junos_services_ssl_initiation_profile" "testacc_security" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_security"
}

resource "junos_security" "testacc_security" {
  alg {
    dns_disable    = true
    ftp_disable    = true
    h323_disable   = true
    mgcp_disable   = true
    msrpc_disable  = true
    pptp_disable   = true
    rsh_disable    = true
    rtsp_disable   = true
    sccp_disable   = true
    sip_disable    = true
    sql_disable    = true
    sunrpc_disable = true
    talk_disable   = true
    tftp_disable   = true
  }
  flow {
    advanced_options {
      drop_matching_reserved_ip_address = true
      drop_matching_link_local_address  = true
      reverse_route_packet_mode_vr      = true
    }
    aging {
      early_ageout   = 10
      high_watermark = 90
      low_watermark  = 80
    }
    allow_dns_reply                       = true
    allow_embedded_icmp                   = true
    allow_reverse_ecmp                    = true
    enable_reroute_uniform_link_check_nat = true
    ethernet_switching {
      block_non_ip_all   = true
      bpdu_vlan_flooding = true
      no_packet_flooding {}
    }
    force_ip_reassembly             = true
    ipsec_performance_acceleration  = true
    mcast_buffer_enhance            = true
    pending_sess_queue_length       = "normal"
    preserve_incoming_fragment_size = true
    route_change_timeout            = 10
    syn_flood_protection_mode       = "syn-proxy"
    sync_icmp_session               = true
    tcp_mss {
      all_tcp_mss = 1499
      gre_in {}
      gre_out {}
      ipsec_vpn {}
    }
    tcp_session {
      fin_invalidate_session = true
      maximum_window         = "512K"
      no_sequence_check      = true
      rst_invalidate_session = true
      rst_sequence_check     = true
      strict_syn_check       = true
      tcp_initial_timeout    = 10
      time_wait_state {}
    }
  }
  forwarding_options {
    inet6_mode            = "flow-based"
    mpls_mode             = "flow-based"
    iso_mode_packet_based = "true"
  }
  forwarding_process {
    enhanced_services_mode = true
  }
  idp_security_package {
    automatic_enable             = true
    automatic_interval           = 24
    automatic_start_time         = "2016-1-1.02:00:00"
    install_ignore_version_check = true
    proxy_profile                = junos_services_proxy_profile.testacc_security.name
    source_address               = "192.0.2.6"
    url                          = "https://signatures.juniper.net/cgi-bin/index.cgi"
  }
  idp_sensor_configuration {
    log_cache_size = 10
    log_suppression {
      disable                        = true
      no_include_destination_address = true
      max_logs_operate               = 1000
      max_time_report                = 30
      start_log                      = 35
    }
    packet_log {
      source_address             = "192.0.2.4"
      host_address               = "192.0.2.5"
      host_port                  = 100
      max_sessions               = 10
      threshold_logging_interval = 20
      total_memory               = 25
    }
    security_configuration_protection_mode = "datacenter"
  }
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      match          = "test"
      size           = 102400
      world_readable = true
    }
    flag            = ["all"]
    rate_limit      = 100
    no_remote_trace = true
  }
  log {
    disable           = true
    facility_override = "local7"
    file {
      files = 10
      name  = "security.log"
      path  = "/"
      size  = 10
    }
    format           = "syslog"
    mode             = "event"
    report           = true
    source_interface = junos_interface_logical.testacc_security.name
    transport {
      protocol        = "tcp"
      tcp_connections = 5
      tls_profile     = junos_services_ssl_initiation_profile.testacc_security.name
    }
    utc_timestamp = true
  }
  nat_source {
    interface_port_overloading_off         = true
    pool_utilization_alarm_raise_threshold = 90
  }
  policies {
    policy_rematch = true
  }
  user_identification_auth_source {
    ad_auth_priority                = 1
    aruba_clearpass_priority        = 2
    firewall_auth_priority          = 3
    local_auth_priority             = 4
    unified_access_control_priority = 0
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
    feature_profile_web_filtering_juniper_enhanced_server {
      host = "192.0.2.1"
      port = 1500
    }
  }
}
