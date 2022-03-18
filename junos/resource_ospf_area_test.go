package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's ge-0/0/4.
func TestAccJunosOspfArea_basic(t *testing.T) {
	testaccOspfArea := defaultInterfaceTestAcc
	testaccOspfArea2 := defaultInterfaceTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccOspfArea = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccOspfArea2 = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
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
			},
		})
	}
}

func testAccJunosOspfAreaConfigCreate(interFace string) string {
	return `
resource junos_ospf_area "testacc_ospfarea" {
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
resource junos_interface_logical "testacc_ospfarea" {
  name        = "` + interFace + `.0"
  description = "testacc_ospfarea"
}
`
}

func testAccJunosOspfAreaConfigUpdate(interFace, interFace2 string) string {
	return `
resource junos_ospf_area "testacc_ospfarea" {
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
}
resource junos_interface_logical "testacc_ospfarea" {
  name        = "` + interFace + `.0"
  description = "testacc_ospfarea"
}
resource junos_interface_logical "testacc_ospfarea2" {
  name             = "` + interFace2 + `.0"
  description      = "testacc_ospfarea2"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
}
resource junos_routing_instance "testacc_ospfarea" {
  name = "testacc_ospfarea"
}
resource junos_ospf_area "testacc_ospfarea2" {
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
`
}
