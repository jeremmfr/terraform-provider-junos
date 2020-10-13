package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosBgpNeighbor_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosBgpNeighborConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"routing_instance", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"group", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_inactive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"as_override", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"damping", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"log_updown", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"multipath", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"remove_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"hold_time", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as", "65001"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as_loops", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"out_delay", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"peer_as", "65002"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"authentication_algorithm", "md5"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_address", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"export.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"export.0", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"import.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"import.0", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.detection_time_threshold", "60"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.transmit_interval_threshold", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.transmit_interval_minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.holddown_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.minimum_receive_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.multiplier", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.0.session_mode", "automatic"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.nlri_type", "unicast"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.0.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.0.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.0.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.0.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.0.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.0.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.nlri_type", "multicast"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.accepted_prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.accepted_prefix_limit.0.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.prefix_limit.0.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet6.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet6.0.accepted_prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet6.0.prefix_limit.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.0.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_bgp_neighbor.testacc_bgpneighbor",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"routing_instance", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"group", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_external_conditional", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"no_advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out_igp_offset", "-10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out_igp_delay_med_update", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"authentication_key", "password"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.0.restart_time", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.0.stale_route_time", "10"),
					),
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate2(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"advertise_external", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"accept_remote_nexthop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"multihop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"local_as_no_prepend_global_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"metric_out_minimum_igp_offset", "-10"),
					),
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate3(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor3",
							"local_as_alias", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor3",
							"metric_out_minimum_igp", "true"),
					),
				},
			},
		})
	}
}

func testAccJunosBgpNeighborConfigCreate() string {
	return `
resource junos_routing_instance "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  as = "65000"
}
resource junos_policyoptions_policy_statement "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  then {
    action = "accept"
  }
}
resource junos_bgp_group "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
}
resource junos_bgp_neighbor "testacc_bgpneighbor" {
  ip = "192.0.2.4"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
  group = junos_bgp_group.testacc_bgpneighbor.name
  advertise_inactive = true
  advertise_peer_as = true
  as_override = true
  damping = true
  log_updown = true
  mtu_discovery = true
  multipath = true
  remove_private = true
  passive = true
  hold_time = 30
  local_as = "65001"
  local_as_private = true
  local_as_loops = 1
  local_preference = 100
  metric_out = 100
  out_delay = 30
  peer_as = "65002"
  preference = 100
  authentication_algorithm = "md5"
  local_address = "192.0.2.3"
  export = [ junos_policyoptions_policy_statement.testacc_bgpneighbor.name ]
  import = [ junos_policyoptions_policy_statement.testacc_bgpneighbor.name ]
  bfd_liveness_detection {
    detection_time_threshold = 60
    transmit_interval_threshold = 30
    transmit_interval_minimum_interval = 10
    holddown_interval = 10
    minimum_interval = 10
    minimum_receive_interval = 10
    multiplier = 2
    session_mode = "automatic"
  }
  family_inet {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum = 2
      teardown = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum = 2
      teardown = 50
      teardown_idle_timeout = 30
    }
  }
  family_inet {
    nlri_type = "multicast"
    accepted_prefix_limit {
      maximum = 2
      teardown_idle_timeout_forever = true
    }
    prefix_limit {
      maximum = 2
      teardown_idle_timeout_forever = true
    }
  }
  family_inet6 {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum = 2
      teardown = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum = 2
      teardown = 50
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
func testAccJunosBgpNeighborConfigUpdate() string {
	return `
resource junos_routing_instance "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  as = "65000"
}
resource junos_policyoptions_policy_statement "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  then {
    action = "accept"
  }
}
resource junos_bgp_group "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
  type = "internal"
}
resource junos_bgp_neighbor "testacc_bgpneighbor" {
  ip = "192.0.2.4"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
  group = junos_bgp_group.testacc_bgpneighbor.name
  advertise_external_conditional = true
  no_advertise_peer_as = true
  metric_out_igp_offset = -10
  metric_out_igp_delay_med_update = true
  authentication_key = "password"
  graceful_restart {
    restart_time = 10
    stale_route_time = 10
  }
}

`
}
func testAccJunosBgpNeighborConfigUpdate2() string {
	return `
resource junos_routing_instance "testacc_bgpneighbor2" {
  name = "testacc_bgpneighbor2"
  as = "65000"
}
resource junos_bgp_group "testacc_bgpneighbor2" {
  name = "testacc_bgpneighbor2"
  routing_instance = junos_routing_instance.testacc_bgpneighbor2.name
  type = "internal"
}
resource junos_bgp_neighbor "testacc_bgpneighbor2" {
  ip = "192.0.2.4"
  routing_instance = junos_routing_instance.testacc_bgpneighbor2.name
  group = junos_bgp_group.testacc_bgpneighbor2.name
  advertise_external = true
  accept_remote_nexthop = true
  multihop = true
  local_as = "65000"
  local_as_no_prepend_global_as = true
  metric_out_minimum_igp_offset = -10
}
`
}
func testAccJunosBgpNeighborConfigUpdate3() string {
	return `
resource junos_routing_instance "testacc_bgpneighbor3" {
  name = "testacc_bgpneighbor3"
  as = "65000"
}
resource junos_bgp_group "testacc_bgpneighbor3" {
  name = "testacc_bgpneighbor3"
  routing_instance = junos_routing_instance.testacc_bgpneighbor3.name
  type = "internal"
}
resource junos_bgp_neighbor "testacc_bgpneighbor3" {
  ip = "192.0.2.4"
  routing_instance = junos_routing_instance.testacc_bgpneighbor3.name
  group = junos_bgp_group.testacc_bgpneighbor3.name
  local_as = "65000"
  local_as_alias = true
  metric_out_minimum_igp = true
}
`
}
