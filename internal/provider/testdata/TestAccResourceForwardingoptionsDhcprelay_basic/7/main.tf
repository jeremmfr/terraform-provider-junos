resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v4_default" {

  no_snoop = true
  overrides_v4 {
    allow_no_end_option             = true
    asymmetric_lease_time           = 7200
    bootp_support                   = true
    delete_binding_on_renegotiation = true
    no_allow_snooped_clients        = true
    no_unicast_replies              = true
  }
  persistent_storage_automatic = true
  relay_option {
    option_60 {
      compare    = "equals"
      value_type = "ascii"
      value      = " equals ascii "
      action     = "drop"
    }
    option_60_default_action {
      action = "forward-only"
    }
    option_77 {
      compare    = "starts-with"
      value_type = "ascii"
      value      = " start ascii "
      action     = "forward-only"
    }
    option_77_default_action {
      action = "drop"
    }
    option_order = ["77", "60"]
  }
  server_response_time = 12001
}

resource "junos_forwardingoptions_dhcprelay" "testacc_dhcprelay_v6_default" {
  version = "v6"

  authentication_username_include {
    client_id                 = true
    relay_agent_interface_id  = true
    relay_agent_remote_id     = true
    relay_agent_subscriber_id = true
  }
  duplicate_clients_incoming_interface = true
  dynamic_profile                      = "junos-default-profile"
  dynamic_profile_use_primary          = "junos-default-dhcp-profile"
  exclude_relay_agent_identifier       = true
  forward_only                         = true
  forward_only_replies                 = true
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
  relay_agent_option_79             = true
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
