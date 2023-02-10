package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's ge-0/0/4.
func TestAccJunosOspfArea_basic(t *testing.T) {
	testaccOspfArea := junos.DefaultInterfaceTestAcc
	testaccOspfArea2 := junos.DefaultInterfaceTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccOspfArea = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccOspfArea2 = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosOspfAreaConfigCreate(testaccOspfArea),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"area_id", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"version", "v2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.name", "all"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.disable", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.passive", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.metric", "100"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.retransmit_interval", "12"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.hello_interval", "11"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.dead_interval", "10"),
					),
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate(testaccOspfArea, testaccOspfArea2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"routing_instance", "testacc_ospfarea"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"version", "v3"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"interface.1.name", testaccOspfArea2+".0"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"interface.1.bfd_liveness_detection.#", "1"),
					),
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfareav3ipv4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea2v3realm",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate2(),
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate3(),
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate4(testaccOspfArea),
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate5(testaccOspfArea),
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea3",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosOspfAreaConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  interface {
    name                = "all"
    disable             = true
    passive             = true
    metric              = 100
    retransmit_interval = 12
    hello_interval      = 11
    dead_interval       = 10
  }
  interface {
    name      = junos_interface_logical.testacc_ospfarea.name
    secondary = true
  }
}
resource "junos_ospf_area" "testacc_ospfareav3ipv4" {
  area_id = "0.0.0.0"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name    = "all"
    disable = true
  }
  interface {
    name      = junos_interface_logical.testacc_ospfarea.name
    secondary = true
  }
}
resource "junos_interface_logical" "testacc_ospfarea" {
  name        = "%s.0"
  description = "testacc_ospfarea"
}
`, interFace)
}

func testAccJunosOspfAreaConfigUpdate(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  interface {
    name                           = "all"
    passive                        = true
    authentication_simple_password = "testPass"
    link_protection                = true
    no_advertise_adjacency_segment = true
    no_interface_state_traps       = true
    no_neighbor_down_notification  = true
    poll_interval                  = 19
    te_metric                      = 221
  }
  interface {
    name    = junos_interface_logical.testacc_ospfarea.name
    disable = true
    authentication_md5 {
      key_id = 3
      key    = "testK3y"
    }
    authentication_md5 {
      key_id     = 2
      key        = "testK3y2"
      start_time = "2022-3-9.12:50:00"
    }
    strict_bfd = true
    bfd_liveness_detection {
      minimum_receive_interval           = 29
      transmit_interval_minimum_interval = 48
      transmit_interval_threshold        = 49
      version                            = "automatic"
    }
    neighbor {
      address = "192.0.2.6"
    }
    neighbor {
      address  = "192.0.2.5"
      eligible = "true"
    }
  }
  network_summary_export = [junos_policyoptions_policy_statement.testacc_ospfarea.name]
  network_summary_import = [junos_policyoptions_policy_statement.testacc_ospfarea2.name]
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfareav3ipv4" {
  area_id = "0.0.0.0"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name     = junos_interface_logical.testacc_ospfarea.name
    priority = 0
    bfd_liveness_detection {
      full_neighbors_only                = true
      minimum_receive_interval           = 27
      transmit_interval_minimum_interval = 50
      transmit_interval_threshold        = 51
    }
  }
}
resource "junos_interface_logical" "testacc_ospfarea" {
  name        = "%s.0"
  description = "testacc_ospfarea"
}
resource "junos_interface_logical" "testacc_ospfarea2" {
  name             = "%s.0"
  description      = "testacc_ospfarea2"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
}
resource "junos_routing_instance" "testacc_ospfarea" {
  name = "testacc_ospfarea"
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id          = "0.0.0.0"
  version          = "v3"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
  interface {
    name                = "all"
    passive             = true
    metric              = 100
    retransmit_interval = 32
    hello_interval      = 31
    dead_interval       = 30
    bandwidth_based_metrics {
      bandwidth = "100k"
      metric    = 13
    }
    bandwidth_based_metrics {
      bandwidth = "1m"
      metric    = 14
    }
    demand_circuit                                    = true
    dynamic_neighbors                                 = true
    flood_reduction                                   = true
    interface_type                                    = "p2mp"
    mtu                                               = 900
    no_eligible_backup                                = true
    no_eligible_remote_backup                         = true
    node_link_protection                              = true
    passive_traffic_engineering_remote_node_id        = "192.0.2.7"
    passive_traffic_engineering_remote_node_router_id = "192.0.2.8"
    priority                                          = 21
    transit_delay                                     = 23
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea2.name
    bfd_liveness_detection {
      authentication_loose_check         = true
      detection_time_threshold           = 60
      full_neighbors_only                = true
      holddown_interval                  = 15
      minimum_interval                   = 16
      minimum_receive_interval           = 17
      multiplier                         = 2
      no_adaptation                      = true
      transmit_interval_minimum_interval = 18
      transmit_interval_threshold        = 19
      version                            = "automatic"
    }
  }
}
resource "junos_ospf_area" "testacc_ospfarea2v3realm" {
  area_id          = "0.0.0.0"
  version          = "v3"
  realm            = "ipv4-multicast"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
  interface {
    name    = "all"
    passive = true
  }
  interface {
    name = junos_interface_logical.testacc_ospfarea2.name
    bfd_liveness_detection {
      version                            = "automatic"
      minimum_receive_interval           = 270
      transmit_interval_minimum_interval = 500
      transmit_interval_threshold        = 510
    }
  }
}
`, interFace, interFace2)
}

func testAccJunosOspfAreaConfigUpdate2() string {
	return `
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  version = "v3"
  interface {
    name    = "all"
    disable = true
  }
  virtual_link {
    neighbor_id  = "192.0.2.0"
    transit_area = "192.0.2.1"
  }
}
`
}

func testAccJunosOspfAreaConfigUpdate3() string {
	return `
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "0.0.0.0"
  version = "v3"
  interface {
    name    = "all"
    disable = true
  }
  virtual_link {
    neighbor_id         = "192.0.2.100"
    transit_area        = "192.0.2.101"
    dead_interval       = 102
    demand_circuit      = true
    disable             = true
    flood_reduction     = true
    hello_interval      = 103
    mtu                 = 1040
    retransmit_interval = 105
    transit_delay       = 106

  }
  virtual_link {
    neighbor_id  = "192.0.2.0"
    transit_area = "192.0.2.1"
  }
}
`
}

func testAccJunosOspfAreaConfigUpdate4(interFace string) string {
	return fmt.Sprintf(`
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "1"
  version = "v3"
  interface {
    name = "all"
  }
  area_range {
    range = "fe80::/64"
  }
  no_context_identifier_advertisement = true
  inter_area_prefix_export = [
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
  ]
  inter_area_prefix_import = [
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
  ]
  nssa {}
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id = "2"
  version = "v3"
  interface {
    name    = "%s.0"
    passive = true
  }
  stub {}
}
`, interFace)
}

func testAccJunosOspfAreaConfigUpdate5(interFace string) string {
	return fmt.Sprintf(`
resource "junos_ospf_area" "testacc_ospfarea" {
  area_id = "1"
  version = "v3"
  interface {
    name = "all"
  }
  area_range {
    range = "fe80:f::/64"
    exact = true
  }
  area_range {
    range           = "fe80:e::/64"
    exact           = true
    override_metric = 106
  }
  area_range {
    range    = "fe80::/64"
    restrict = true
  }
  context_identifier = ["127.0.0.2", "127.0.0.1"]
  inter_area_prefix_export = [
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
  ]
  inter_area_prefix_import = [
    junos_policyoptions_policy_statement.testacc_ospfarea2.name,
    junos_policyoptions_policy_statement.testacc_ospfarea.name,
  ]
  nssa {
    area_range {
      range = "fe80::/64"
      exact = true
    }
    area_range {
      range           = "fe80:b::/64"
      override_metric = 107
    }
    area_range {
      range    = "fe80:a::/64"
      restrict = true
    }
    default_lsa {
      default_metric = 109
      metric_type    = 2
      type_7         = true
    }
    summaries = true
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea" {
  name = "testacc_ospfarea"
  then {
    action = "accept"
  }
}
resource "junos_policyoptions_policy_statement" "testacc_ospfarea2" {
  name = "testacc_ospfarea2"
  then {
    action = "reject"
  }
}
resource "junos_ospf_area" "testacc_ospfarea2" {
  area_id = "2"
  version = "v3"
  interface {
    name    = "%s.0"
    passive = true
  }
  stub {
    default_metric = 150
    no_summaries   = true
  }
}
resource "junos_ospf_area" "testacc_ospfarea3" {
  area_id = "3"
  version = "v3"
  realm   = "ipv4-unicast"
  interface {
    name    = "%s.0"
    passive = true
  }
  nssa {
    no_summaries = true
    default_lsa {}
  }
}
`, interFace, interFace)
}
