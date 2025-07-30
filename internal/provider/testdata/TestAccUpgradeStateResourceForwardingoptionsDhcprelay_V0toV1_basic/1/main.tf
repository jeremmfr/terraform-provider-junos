resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v4_default" {
  active_leasequery {
    idle_timeout      = 10
    peer_address      = "192.0.2.1"
    timeout           = 11
    topology_discover = true
  }
  active_server_group                     = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v4_default.name
  active_server_group_allow_server_change = true
  arp_inspection                          = true
  authentication_password                 = "test1#1"
  authentication_username_include {
    circuit_type                               = true
    client_id                                  = true
    client_id_exclude_headers                  = true
    client_id_use_automatic_ascii_hex_encoding = true
    delimiter                                  = "#"
    domain_name                                = "a domain"
    interface_description                      = "logical"
    interface_name                             = true
    mac_address                                = true
    option_60                                  = true
    option_82                                  = true
    option_82_circuit_id                       = true
    option_82_remote_id                        = true
    routing_instance_name                      = true
    user_prefix                                = "user_#1_"
    vlan_tags                                  = true
  }
  bulk_leasequery {}
  client_response_ttl           = 60
  duplicate_clients_in_subnet   = "option-82"
  dynamic_profile               = "junos-default-profile"
  forward_only                  = true
  forward_only_routing_instance = junos_routing_instance.testacc_dhcprelay.name
  maximum_hop_count             = 8
  minimum_wait_time             = 0

  overrides_v4 {
    allow_snooped_clients         = true
    always_write_giaddr           = true
    always_write_option_82        = true
    delay_authentication          = true
    disable_relay                 = true
    layer2_unicast_replies        = true
    no_bind_on_request            = true
    proxy_mode                    = true
    relay_source                  = "lo0.1"
    replace_ip_source_with_giaddr = true
    send_release_on_delete        = true
    trust_option_82               = true
    user_defined_option_82        = "#test"
  }
  relay_option {
    option_60 {
      compare    = "equals"
      value_type = "ascii"
      value      = " equals ascii "
      action     = "local-server-group"
      group      = junos_system_services_dhcp_localserver_group.testacc_dhcprelay_v4_default.name
    }
    option_60_default_action {
      action = "relay-server-group"
      group  = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v4_default.name
    }
    option_77 {
      compare    = "starts-with"
      value_type = "ascii"
      value      = " start ascii "
      action     = "relay-server-group"
      group      = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v4_default.name
    }
    option_77 {
      compare    = "equals"
      value_type = "hexadecimal"
      value      = "11BBee"
      action     = "relay-server-group"
      group      = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v4_default.name
    }
    option_77_default_action {
      action = "local-server-group"
      group  = junos_system_services_dhcp_localserver_group.testacc_dhcprelay_v4_default.name
    }
    option_order = ["77", "60"]
  }
  relay_option_82 {
    circuit_id {
      include_irb_and_l2           = true
      keep_incoming_circuit_id     = true
      no_vlan_interface_name       = true
      prefix_host_name             = true
      prefix_routing_instance_name = true
      use_interface_description    = "logical"
    }
    exclude_relay_agent_identifier = true
    link_selection                 = true
    remote_id {
      include_irb_and_l2           = true
      keep_incoming_remote_id      = true
      no_vlan_interface_name       = true
      prefix_routing_instance_name = true
      use_interface_description    = "device"
    }
    vendor_specific_host_name = true
    vendor_specific_location  = true
  }
  route_suppression_destination = true
  server_match_address {
    address = "192.0.2.1/30"
    action  = "create-relay-entry"
  }
  server_match_default_action = "forward-only"
  source_ip_change            = true
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_v4_default" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_dhcprelay_v4_default"
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcprelay_v4_default" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_dhcprelay_v4_default"

  dynamic_profile = "junos-default-profile"
}

resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v6_default" {
  version = "v6"

  authentication_username_include {
    client_id                 = true
    relay_agent_interface_id  = true
    relay_agent_remote_id     = true
    relay_agent_subscriber_id = true
  }
  bulk_leasequery {
    attempts          = 4
    timeout           = 2
    trigger_automatic = true
  }
  duplicate_clients_incoming_interface = true
  dynamic_profile                      = "junos-default-profile"
  dynamic_profile_use_primary          = "junos-default-dhcp-profile"
  exclude_relay_agent_identifier       = true
  forward_only                         = true
  forward_only_replies                 = true
  forward_only_routing_instance        = junos_routing_instance.testacc_dhcprelay.name
  lease_time_validation {}
  liveness_detection_method_layer2 {
    max_consecutive_retries = 4
    transmit_interval       = 305
  }
  overrides_v6 {
    always_process_option_request_option        = true
    asymmetric_lease_time                       = 900
    asymmetric_prefix_lease_time                = 1000
    client_negotiation_match_incoming_interface = true
    delay_authentication                        = true
    delete_binding_on_renegotiation             = true
    dual_stack                                  = "dual-#stack"
    interface_client_limit                      = 120
    no_allow_snooped_clients                    = true
    no_bind_on_request                          = true
    relay_source                                = "lo0.1"
    send_release_on_delete                      = true
  }
  relay_agent_interface_id {}
  relay_agent_option_79 = true
  relay_agent_remote_id {}
  relay_option {
    option_15 {
      compare    = "equals"
      value_type = "ascii"
      value      = "&equals-ascii "
      action     = "drop"
    }
    option_15 {
      compare    = "equals"
      value_type = "hexadecimal"
      value      = "AABBff"
      action     = "relay-server-group"
      group      = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v6_default.name
    }
    option_15_default_action {
      action = "relay-server-group"
      group  = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v6_default.name
    }
    option_16 {
      compare    = "starts-with"
      value_type = "ascii"
      value      = " start ascii "
      action     = "relay-server-group"
      group      = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v6_default.name
    }
    option_16 {
      compare    = "equals"
      value_type = "hexadecimal"
      value      = "11BBee"
      action     = "relay-server-group"
      group      = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v6_default.name
    }
    option_16_default_action {
      action = "relay-server-group"
      group  = junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_v6_default.name
    }
    option_order = ["16", "15"]
  }
  route_suppression_access          = true
  route_suppression_access_internal = true
  server_match_duid {
    compare    = "equals"
    value_type = "ascii"
    value      = " test_space "
    action     = "forward-only"
  }
  vendor_specific_information_host_name = true
  vendor_specific_information_location  = true
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_v6_default" {
  lifecycle {
    create_before_destroy = true
  }
  name    = "testacc_dhcprelay_v6_default"
  version = "v6"
}

resource "junos_routing_instance" "testacc_dhcprelay" {
  name = "testacc_dhcprelay"
}
resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v4_ri" {
  routing_instance = junos_routing_instance.testacc_dhcprelay.name

  authentication_username_include {
    client_id                                  = true
    client_id_exclude_headers                  = true
    client_id_use_automatic_ascii_hex_encoding = true
    option_82                                  = true
  }
  dynamic_profile                   = "junos-default-profile"
  dynamic_profile_aggregate_clients = true
  forward_snooped_clients           = "all-interfaces"
  lease_time_validation {
    lease_time_threshold  = 60099
    violation_action_drop = true
  }
  liveness_detection_method_bfd {
    detection_time_threshold    = 200000
    holddown_interval           = 2
    minimum_interval            = 30003
    minimum_receive_interval    = 30004
    multiplier                  = 5
    no_adaptation               = true
    session_mode                = "multihop"
    transmit_interval_minimum   = 30006
    transmit_interval_threshold = 30066
    version                     = "automatic"
  }
  liveness_detection_failure_action = "log-only"
  leasequery {
    attempts = 9
    timeout  = 8
  }
  overrides_v4 {
    client_discover_match  = "incoming-interface"
    always_write_option_82 = true
  }
  relay_option_82 {
    circuit_id {
      use_vlan_id = true
    }
    remote_id {
      use_vlan_id = true
    }
  }
  route_suppression_destination           = true
  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}

resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v6_ri" {
  routing_instance = junos_routing_instance.testacc_dhcprelay.name
  version          = "v6"


  dynamic_profile                          = "junos-default-profile"
  dynamic_profile_aggregate_clients        = true
  dynamic_profile_aggregate_clients_action = "merge"

  overrides_v6 {
    allow_snooped_clients = true
  }
  relay_agent_interface_id {
    include_irb_and_l2           = true
    keep_incoming_id             = true
    keep_incoming_id_strict      = true
    no_vlan_interface_name       = true
    prefix_host_name             = true
    prefix_routing_instance_name = true
    use_interface_description    = "logical"
    use_option_82                = true
  }
  relay_agent_remote_id {
    keep_incoming_id     = true
    use_option_82        = true
    use_option_82_strict = true
    use_vlan_id          = true
  }
  relay_option {
    option_15 {
      compare    = "equals"
      value_type = "ascii"
      value      = " equals ascii "
      action     = "forward-only"
    }
    option_15 {
      compare    = "equals"
      value_type = "hexadecimal"
      value      = "11BBee"
      action     = "drop"
    }
    option_15_default_action {
      action = "drop"
    }
    option_16 {
      compare    = "starts-with"
      value_type = "ascii"
      value      = " start ascii "
      action     = "forward-only"
    }
    option_16_default_action {
      action = "forward-only"
    }
    option_order = ["15", "16"]
  }
  remote_id_mismatch_disconnect           = true
  route_suppression_access                = true
  service_profile                         = "service-pro#2"
  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
