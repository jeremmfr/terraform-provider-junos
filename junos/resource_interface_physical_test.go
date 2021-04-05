package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE_AE=ae<num> for choose interface aggregate test else it's ae0.
func TestAccJunosInterfacePhysical_basic(t *testing.T) {
	var testaccInterface string
	var testaccInterface2 string
	var testaccInterfaceAE string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	if os.Getenv("TESTACC_INTERFACE_AE") != "" {
		testaccInterfaceAE = os.Getenv("TESTACC_INTERFACE_AE")
	} else {
		testaccInterfaceAE = "ae0"
	}
	if os.Getenv("TESTACC_INTERFACE2") != "" {
		testaccInterface2 = os.Getenv("TESTACC_INTERFACE2")
	} else {
		testaccInterface2 = "ge-0/0/4"
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfacePhysicalSWConfigCreate(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interface"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"trunk", "true"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_native", "100"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.0", "100-110"),
					),
				},
				{
					Config: testAccJunosInterfacePhysicalSWConfigUpdate(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interfaceU"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"trunk", "false"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_native", "0"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.0", "100"),
					),
				},
				{
					ResourceName:      "junos_interface_physical.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		if os.Getenv("TESTACC_ROUTER") != "" {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccJunosInterfacePhysicalRouterConfigCreate(testaccInterface, testaccInterfaceAE),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"parent_ether_opts.#", "1"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"parent_ether_opts.0.source_address_filter.#", "1"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"parent_ether_opts.0.source_filtering", "true"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.#", "1"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.0.identifier", "00:01:11:11:11:11:11:11:11:11"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.0.mode", "all-active"),
						),
					},
					{
						Config: testAccJunosInterfacePhysicalRouterConfigUpdate(testaccInterface, testaccInterfaceAE),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.#", "1"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.0.identifier", "00:11:11:11:11:11:11:11:11:11"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
								"esi.0.mode", "all-active"),
						),
					},
					{
						ResourceName:      "junos_interface_physical.testacc_interfaceAE",
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		}
		if os.Getenv("TESTACC_SRX") != "" {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccJunosInterfacePhysicalSRXConfigCreate(testaccInterface, testaccInterface2),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface_reth",
								"parent_ether_opts.#", "1"),
							resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface_reth",
								"parent_ether_opts.0.redundancy_group", "1"),
						),
					},
				},
			})
		}
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfacePhysicalConfigCreate(testaccInterface, testaccInterfaceAE, testaccInterface2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interface"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"gigether_opts.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"gigether_opts.0.ae_8023ad", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"name", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.0.lacp.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.0.lacp.0.mode", "active"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.0.minimum_links", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"vlan_tagging", "true"),
					),
				},
				{
					Config: testAccJunosInterfacePhysicalConfigUpdate(testaccInterface, testaccInterfaceAE, testaccInterface2),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interfaceU"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.0.lacp.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.0.minimum_bandwidth", "1 gbps"),
					),
				},
				{
					ResourceName:      "junos_interface_physical.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_interface_physical.testacc_interfaceAE",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosInterfacePhysicalConfigUpdate2(testaccInterface, testaccInterfaceAE),
				},
			},
		})
	}
}

func testAccJunosInterfacePhysicalSWConfigCreate(interFace string) string {
	return `
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interface"
  trunk        = true
  vlan_native  = 100
  vlan_members = ["100-110"]
}
`
}

func testAccJunosInterfacePhysicalSWConfigUpdate(interFace string) string {
	return `
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interfaceU"
  vlan_members = ["100"]
}
`
}

func testAccJunosInterfacePhysicalConfigCreate(interFace, interfaceAE, interFace2 string) string {
	return `
resource junos_interface_physical testacc_interface {
  name        = "` + interFace + `"
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = "` + interfaceAE + `"
  }
}
resource junos_interface_physical testacc_interface2 {
  name        = "` + interFace2 + `"
  description = "testacc_interface2"
  gigether_opts {
    flow_control     = true
    loopback         = true
    auto_negotiation = true
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [
    junos_interface_physical.testacc_interface,
  ]
  name        = "` + interfaceAE + `"
  description = "testacc_interfaceAE"
  parent_ether_opts {
    flow_control = true
    lacp {
      mode            = "active"
      admin_key       = 1
      periodic        = "slow"
      sync_reset      = "disable"
      system_id       = "00:00:01:00:01:00"
      system_priority = 250
    }
    loopback      = true
    link_speed    = "1g"
    minimum_links = 1
  }
  vlan_tagging = true
}
`
}

func testAccJunosInterfacePhysicalConfigUpdate(interFace, interfaceAE, interFace2 string) string {
	return `
resource junos_interface_physical testacc_interface {
  name        = "` + interFace + `"
  description = "testacc_interfaceU"
  gigether_opts {
    ae_8023ad = "` + interfaceAE + `"
  }
}
resource junos_interface_physical testacc_interface2 {
  name        = "` + interFace2 + `"
  description = "testacc_interface2"
  ether_opts {
    flow_control     = true
    loopback         = true
    auto_negotiation = true
  }
}
resource "junos_interface_logical" "testacc_interfaceLO" {
  name = "lo0.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/32"
    }
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  depends_on = [
    junos_interface_physical.testacc_interface,
    junos_interface_logical.testacc_interfaceLO,
  ]
  name        = "` + interfaceAE + `"
  description = "testacc_interfaceAE"
  parent_ether_opts {
    bfd_liveness_detection {
      local_address                      = "192.0.2.1"
      detection_time_threshold           = 30
      holddown_interval                  = 30
      minimum_interval                   = 30
      minimum_receive_interval           = 10
      multiplier                         = 1
      neighbor                           = "192.0.2.2"
      no_adaptation                      = true
      transmit_interval_minimum_interval = 10
      transmit_interval_threshold        = 30
      version                            = "automatic"
    }
    no_flow_control   = true
    no_loopback       = true
    link_speed        = "1g"
    minimum_bandwidth = "1 gbps"
  }
  vlan_tagging = true
}
`
}

func testAccJunosInterfacePhysicalConfigUpdate2(interFace, interfaceAE string) string {
	return `
resource junos_interface_physical testacc_interface {
  name        = "` + interFace + `"
  description = "testacc_interfaceU"
  ether_opts {
    ae_8023ad = "` + interfaceAE + `"
  }
}
`
}

func testAccJunosInterfacePhysicalRouterConfigCreate(interFace, interfaceAE string) string {
	return `
resource "junos_interface_physical" "testacc_interface" {
  name        = "` + interFace + `"
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = "` + interfaceAE + `"
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  name        = "` + interfaceAE + `"
  description = "testacc_interfaceAE"
  esi {
    identifier = "00:01:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  parent_ether_opts {
    source_address_filter = ["00:11:22:33:44:55"]
    source_filtering      = true
  }
  vlan_tagging = true
}
`
}

func testAccJunosInterfacePhysicalRouterConfigUpdate(interFace, interfaceAE string) string {
	return `
resource "junos_interface_physical" "testacc_interface" {
  name        = "` + interFace + `"
  description = "testacc_interface"
  gigether_opts {
    ae_8023ad = "` + interfaceAE + `"
  }
}
resource "junos_interface_physical" "testacc_interfaceAE" {
  name        = "` + interfaceAE + `"
  description = "testacc_interfaceAE"
  esi {
    identifier = "00:11:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
  vlan_tagging = true
}
`
}

func testAccJunosInterfacePhysicalSRXConfigCreate(interFace, interFace2 string) string {
	return `
resource "junos_interface_physical" "testacc_interface" {
  depends_on = [
    junos_chassis_cluster.testacc_interface
  ]
  name        = "` + interFace + `"
  description = "testacc_interface"
  gigether_opts {
    redundant_parent = "reth0"
  }
}
resource "junos_chassis_cluster" "testacc_interface" {
  fab0 {
    member_interfaces = ["` + interFace2 + `"]
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  redundancy_group {
    node0_priority = 100
    node1_priority = 99
  }
  reth_count = 1
}
resource "junos_interface_physical" "testacc_interface_reth" {
  depends_on = [
    junos_interface_physical.testacc_interface
  ]
  name        = "reth0"
  description = "testacc_interface_reth"
  parent_ether_opts {
    redundancy_group = 1
    minimum_links    = 1
  }
}
`
}
