package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSystemServicesDhcpLocalserverGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSystemServicesDhcpLocalserverGroupConfigCreate(),
				},
				{
					Config: testAccResourceSystemServicesDhcpLocalserverGroupConfigUpdate(),
				},
				{
					ResourceName:      "junos_system_services_dhcp_localserver_group.testacc_dhcpgroup_v4_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_system_services_dhcp_localserver_group.testacc_dhcpgroup_v6_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_system_services_dhcp_localserver_group.testacc_dhcpgroup_v4_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_system_services_dhcp_localserver_group.testacc_dhcpgroup_v6_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceSystemServicesDhcpLocalserverGroupConfigCreate() string {
	return `
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_default" {
  name = "testacc_dhcpgroup_v4_default"

  dynamic_profile = "junos-default-profile"
}

resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_default" {
  name    = "testacc_dhcpgroup_v6_default"
  version = "v6"

  interface {
    name    = "ge-0/0/3.1"
    upto    = "ge-0/0/3.3"
    exclude = true
  }
}

resource "junos_routing_instance" "testacc_dhcpgroup" {
  name = "testacc_dhcpgroup"
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_ri" {
  name             = "testacc_dhcpgroup_v4_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name

  interface {
    name            = "ge-0/0/3.0"
    dynamic_profile = "junos-default-profile"
    trace           = true
  }
  interface {
    name            = "ge-0/0/3.1"
    dynamic_profile = "junos-default-profile"
  }
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_ri" {
  name             = "testacc_dhcpgroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name
  version          = "v6"

  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
`
}

func testAccResourceSystemServicesDhcpLocalserverGroupConfigUpdate() string {
	return `
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_default" {
  name = "testacc_dhcpgroup_v4_default"

  dynamic_profile         = "junos-default-profile"
  authentication_password = "test1#1"
  authentication_username_include {
    circuit_type          = true
    delimiter             = "#"
    domain_name           = "a domain"
    interface_description = "logical"
    interface_name        = true
    mac_address           = true
    option_60             = true
    option_82             = true
    option_82_circuit_id  = true
    option_82_remote_id   = true
    routing_instance_name = true
    user_prefix           = "user_#1_"
    vlan_tags             = true
  }
  reauthenticate_lease_renewal            = true
  reauthenticate_remote_id_mismatch       = true
  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}

resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_default" {
  name    = "testacc_dhcpgroup_v6_default"
  version = "v6"

  authentication_username_include {
    client_id                 = true
    relay_agent_interface_id  = true
    relay_agent_remote_id     = true
    relay_agent_subscriber_id = true
  }
  dynamic_profile             = "junos-default-profile"
  dynamic_profile_use_primary = "junos-default-dhcp-profile"
  interface {
    name    = "ge-0/0/3.1"
    upto    = "ge-0/0/3.3"
    exclude = true
  }
  lease_time_validation {}
  liveness_detection_method_layer2 {
    max_consecutive_retries = 4
    transmit_interval       = 305
  }
  overrides_v6 {
    always_add_option_dns_server                = true
    always_process_option_request_option        = true
    asymmetric_lease_time                       = 900
    asymmetric_prefix_lease_time                = 1000
    client_negotiation_match_incoming_interface = true
    delay_advertise_based_on {
      option     = "option-15"
      compare    = "equals"
      value_type = "ascii"
      value      = "1#1 2"
    }
    delay_advertise_delay_time             = 11
    delegated_pool                         = "test foo"
    delete_binding_on_renegotiation        = true
    dual_stack                             = "dual-#stack"
    interface_client_limit                 = 120
    multi_address_embedded_option_response = true
    process_inform                         = true
    protocol_attributes                    = "test#proto"
    rapid_commit                           = true
    top_level_status_code                  = true
  }
  route_suppression_access_internal = true
}

resource "junos_routing_instance" "testacc_dhcpgroup" {
  name = "testacc_dhcpgroup"
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v4_ri" {
  name             = "testacc_dhcpgroup_v4_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name

  authentication_username_include {
    client_id                                  = true
    client_id_exclude_headers                  = true
    client_id_use_automatic_ascii_hex_encoding = true
    option_82                                  = true
  }
  dynamic_profile                   = "junos-default-profile"
  dynamic_profile_aggregate_clients = true
  interface {
    name                        = "ge-0/0/3.0"
    dynamic_profile             = "junos-default-profile"
    dynamic_profile_use_primary = "junos-default-profile"
    trace                       = true
  }
  interface {
    name            = "ge-0/0/3.1"
    dynamic_profile = "junos-default-profile"
  }
  interface {
    name                                     = "ge-0/0/3.4"
    dynamic_profile                          = "junos-default-profile"
    dynamic_profile_aggregate_clients        = true
    dynamic_profile_aggregate_clients_action = "merge"
    overrides_v4 {
      allow_no_end_option             = true
      asymmetric_lease_time           = 900
      bootp_support                   = true
      client_discover_match           = "option60-and-option82"
      delete_binding_on_renegotiation = true
      dual_stack                      = "dual-#stack"
      include_option_82_forcerenew    = true
      include_option_82_nak           = true
      interface_client_limit          = 120
      process_inform                  = true
      protocol_attributes             = "test#proto"
    }
    service_profile                         = "a_service#1"
    short_cycle_protection_lockout_max_time = 2
    short_cycle_protection_lockout_min_time = 1
  }
  lease_time_validation {
    lease_time_threshold = 60099
    violation_action     = "strict"
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
  overrides_v4 {
    client_discover_match = "incoming-interface"
    delay_offer_based_on {
      option     = "option-60"
      compare    = "equals"
      value_type = "ascii"
      value      = "test #1"
    }
    delay_offer_delay_time = 2
    process_inform         = true
    process_inform_pool    = "pool@21#"
  }
  reconfigure {}
  route_suppression_destination           = true
  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
resource "junos_system_services_dhcp_localserver_group" "testacc_dhcpgroup_v6_ri" {
  name             = "testacc_dhcpgroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcpgroup.name
  version          = "v6"


  dynamic_profile                          = "junos-default-profile"
  dynamic_profile_aggregate_clients        = true
  dynamic_profile_aggregate_clients_action = "merge"

  reconfigure {
    attempts                  = 1
    clear_on_abort            = true
    support_option_pd_exclude = true
    timeout                   = 2
    token                     = "tok en #"
    trigger_radius_disconnect = true
  }
  remote_id_mismatch_disconnect           = true
  route_suppression_access                = true
  service_profile                         = "service-pro#2"
  short_cycle_protection_lockout_max_time = 2
  short_cycle_protection_lockout_min_time = 1
}
`
}
