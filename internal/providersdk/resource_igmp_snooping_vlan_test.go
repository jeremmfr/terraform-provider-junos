package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccJunosIgmpSnoopingVlan_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := defaultInterfaceSwitchTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosIgmpSnoopingVlanSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosIgmpSnoopingVlanSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_igmp_snooping_vlan.vlan10",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		testaccInterface := defaultInterfaceTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosIgmpSnoopingVlanConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosIgmpSnoopingVlanConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_igmp_snooping_vlan.vlan10",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosIgmpSnoopingVlanSWConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_igmp_snooping_vlan" "all" {
  name = "all"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name            = "vlan10"
  immediate_leave = true
  interface {
    name                = "%s.1"
    host_only_interface = true
  }
  proxy = true
}
`, interFace)
}

func testAccJunosIgmpSnoopingVlanSWConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_igmp_snooping_vlan" "vlan10" {
  name            = "vlan10"
  immediate_leave = true
  interface {
    name = "ge-0/0/3.1"
  }
  interface {
    name                       = "%s.0"
    group_limit                = 32
    immediate_leave            = true
    multicast_router_interface = true
    static_group {
      address = "224.255.0.2"
    }
    static_group {
      address = "224.255.0.1"
      source  = "192.0.2.1"
    }
  }
  l2_querier_source_address  = "192.0.2.10"
  proxy                      = true
  proxy_source_address       = "192.0.2.11"
  query_interval             = 33
  query_last_member_interval = "32.1"
  query_response_interval    = "31.0"
  robust_count               = 5
}
`, interFace)
}

func testAccJunosIgmpSnoopingVlanConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_routing_instance" "testacc_igmp_snooping_vlan" {
  name = "testacc_igmp_snooping_vlan"
  type = "vpls"
}
resource "junos_igmp_snooping_vlan" "all" {
  name = "all"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name             = "vlan10"
  routing_instance = junos_routing_instance.testacc_igmp_snooping_vlan.name
  immediate_leave  = true
  interface {
    name                = "%s.1"
    host_only_interface = true
  }
  proxy = true
}
`, interFace)
}

func testAccJunosIgmpSnoopingVlanConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_routing_instance" "testacc_igmp_snooping_vlan" {
  name = "testacc_igmp_snooping_vlan"
  type = "vpls"
}
resource "junos_igmp_snooping_vlan" "vlan10" {
  name             = "vlan10"
  routing_instance = junos_routing_instance.testacc_igmp_snooping_vlan.name
  immediate_leave  = true
  interface {
    name = "ge-0/0/3.1"
  }
  interface {
    name                       = "%s.0"
    group_limit                = 32
    immediate_leave            = true
    multicast_router_interface = true
    static_group {
      address = "224.255.0.2"
    }
    static_group {
      address = "224.255.0.1"
      source  = "192.0.2.1"
    }
  }
  proxy                      = true
  proxy_source_address       = "192.0.2.11"
  query_interval             = 33
  query_last_member_interval = "32.1"
  query_response_interval    = "31.0"
  robust_count               = 5
}
`, interFace)
}
