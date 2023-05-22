package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosBgpGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosBgpGroupConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"routing_instance", "testacc_bgpgroup"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"type", "external"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"advertise_inactive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"as_override", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"cluster", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"damping", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"log_updown", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"remove_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"hold_time", "30"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_as", "65001"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_as_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_as_loops", "1"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"metric_out", "100"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"out_delay", "30"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"peer_as", "65002"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"authentication_algorithm", "md5"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_address", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"export.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"export.0", "testacc_bgpgroup"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"import.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"import.0", "testacc_bgpgroup"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.detection_time_threshold", "60"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.transmit_interval_threshold", "30"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.transmit_interval_minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.holddown_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.minimum_receive_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.multiplier", "2"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bfd_liveness_detection.session_mode", "automatic"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.nlri_type", "unicast"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.accepted_prefix_limit.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.accepted_prefix_limit.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.accepted_prefix_limit.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.prefix_limit.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.prefix_limit.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.0.prefix_limit.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.1.nlri_type", "multicast"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.1.accepted_prefix_limit.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet.1.prefix_limit.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"family_inet6.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"graceful_restart.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_bgp_group.testacc_bgpgroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosBgpGroupConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"routing_instance", "testacc_bgpgroup"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"advertise_external_conditional", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"bgp_multipath.multiple_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"no_advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"metric_out_igp_offset", "-10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"metric_out_igp_delay_med_update", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"authentication_key", "password"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"graceful_restart.restart_time", "10"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"graceful_restart.stale_route_time", "10"),
					),
				},
				{
					Config: testAccJunosBgpGroupConfigUpdate2(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"type", "internal"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"advertise_external", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"accept_remote_nexthop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"multihop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_as_no_prepend_global_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"metric_out_minimum_igp_offset", "-10"),
					),
				},
				{
					Config: testAccJunosBgpGroupConfigUpdate3(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"local_as_alias", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"metric_out_minimum_igp", "true"),
						resource.TestCheckResourceAttr("junos_bgp_group.testacc_bgpgroup",
							"type", "external"),
					),
				},
			},
		})
	}
}

func testAccJunosBgpGroupConfigCreate() string {
	return `
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
`
}

func testAccJunosBgpGroupConfigUpdate() string {
	return `
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
resource "junos_bgp_group" "testacc_bgpgroup" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup
  ]
  name                            = "testacc_bgpgroup"
  description                     = "testacc bgpgroup"
  routing_instance                = junos_routing_instance.testacc_bgpgroup.name
  advertise_external_conditional  = true
  keep_all                        = true
  no_advertise_peer_as            = true
  metric_out_igp_offset           = -10
  metric_out_igp_delay_med_update = true
  authentication_key              = "password"
  bgp_multipath {
    multiple_as = true
  }
  graceful_restart {
    restart_time     = 10
    stale_route_time = 10
  }
  tcp_aggressive_transmission = true
  bgp_error_tolerance {}
}
`
}

func testAccJunosBgpGroupConfigUpdate2() string {
	return `
resource "junos_routing_options" "testacc_bgpgroup" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_bgp_group" "testacc_bgpgroup" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup
  ]
  name                          = "testacc_bgpgroup"
  advertise_external            = true
  accept_remote_nexthop         = true
  multihop                      = true
  local_as                      = "65000"
  local_as_no_prepend_global_as = true
  metric_out_minimum_igp_offset = -10
  type                          = "internal"
  family_evpn {
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
  bgp_error_tolerance {
    malformed_route_limit         = 234
    malformed_update_log_interval = 567
  }
}
`
}

func testAccJunosBgpGroupConfigUpdate3() string {
	return `
resource "junos_routing_options" "testacc_bgpgroup" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_bgp_group" "testacc_bgpgroup" {
  depends_on = [
    junos_routing_options.testacc_bgpgroup
  ]
  name                   = "testacc_bgpgroup"
  local_as               = "65000"
  local_as_alias         = true
  metric_out_minimum_igp = true
  family_evpn {}
  bgp_error_tolerance {
    no_malformed_route_limit = true
  }
}
`
}
