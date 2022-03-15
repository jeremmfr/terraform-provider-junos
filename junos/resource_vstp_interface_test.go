package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's xe-0/0/3.
func TestAccJunosVstpInterface_basic(t *testing.T) {
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
					Config: testAccJunosVstpInterfaceSWConfigCreate(testaccInterface),
				},
				{
					Config: testAccJunosVstpInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface3",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface5",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_vstp_interface.testacc_vstp_interface6",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosVstpInterfaceSWConfigCreate(interFace string) string {
	return `
resource "junos_vstp_interface" "testacc_vstp_interface" {
  name = "all"
}
resource "junos_vstp_vlan" "testacc_vstp_interface2" {
  vlan_id = "10"
}
resource junos_interface_physical testacc_vstp_interface2 {
  name         = "` + interFace + `"
  description  = "testacc_vstp_interface2"
  vlan_members = ["15"]
}
resource "junos_vstp_interface" "testacc_vstp_interface2" {
  name = junos_interface_physical.testacc_vstp_interface2.name
  vlan = junos_vstp_vlan.testacc_vstp_interface2.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface3" {
  name = "testacc_vstp_interface2"
  vlan = ["11"]
}
resource "junos_vstp_interface" "testacc_vstp_interface3" {
  name       = "all"
  vlan_group = junos_vstp_vlan_group.testacc_vstp_interface3.name
}
resource "junos_routing_instance" "testacc_vstp_interface" {
  name = "testacc_vstp_intface"
  type = "virtual-switch"
}
resource "junos_vstp_interface" "testacc_vstp_interface4" {
  name             = "all"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_vlan" "testacc_vstp_interface5" {
  vlan_id          = "all"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_interface" "testacc_vstp_interface5" {
  name             = "all"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = junos_vstp_vlan.testacc_vstp_interface5.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface6" {
  name             = "testacc_vstp_interface6"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = ["13"]
}
resource "junos_vstp_interface" "testacc_vstp_interface6" {
  name             = "` + interFace + `"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan_group       = junos_vstp_vlan_group.testacc_vstp_interface6.name
}
`
}

func testAccJunosVstpInterfaceSWConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_vstp_interface" "testacc_vstp_interface" {
  name                      = "all"
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 101
  edge                      = true
  priority                  = 32
}
resource "junos_vstp_vlan" "testacc_vstp_interface2" {
  vlan_id = "10"
}
resource junos_interface_physical testacc_vstp_interface2 {
  name         = "` + interFace + `"
  description  = "testacc_vstp_interface2"
  vlan_members = ["15"]
}
resource "junos_vstp_interface" "testacc_vstp_interface2" {
  name         = junos_interface_physical.testacc_vstp_interface2.name
  access_trunk = true
  mode         = "shared"
  no_root_port = true
  vlan         = junos_vstp_vlan.testacc_vstp_interface2.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface3" {
  name = "testacc_vstp_interface2"
  vlan = ["11"]
}
resource "junos_vstp_interface" "testacc_vstp_interface3" {
  name       = "all"
  priority   = 32
  vlan_group = junos_vstp_vlan_group.testacc_vstp_interface3.name
}
resource "junos_routing_instance" "testacc_vstp_interface" {
  name = "testacc_vstp_intface"
  type = "virtual-switch"
}
resource "junos_vstp_interface" "testacc_vstp_interface4" {
  name             = "all"
  edge             = true
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_vlan" "testacc_vstp_interface5" {
  vlan_id          = "all"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
}
resource "junos_vstp_interface" "testacc_vstp_interface5" {
  name             = "all"
  mode             = "point-to-point"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = junos_vstp_vlan.testacc_vstp_interface5.vlan_id
}
resource "junos_vstp_vlan_group" "testacc_vstp_interface6" {
  name             = "testacc_vstp_interface6"
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan             = ["13"]
}
resource "junos_vstp_interface" "testacc_vstp_interface6" {
  name             = "` + interFace + `"
  no_root_port     = true
  priority         = 64
  routing_instance = junos_routing_instance.testacc_vstp_interface.name
  vlan_group       = junos_vstp_vlan_group.testacc_vstp_interface6.name
}
`)
}
