package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's xe-0/0/3.
func TestAccJunosRstpInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface := junos.DefaultInterfaceSwitchTestAcc
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRstpInterfaceSWConfigCreate(),
				},
				{
					Config: testAccJunosRstpInterfaceSWConfigUpdate(testaccInterface),
				},
				{
					ResourceName:      "junos_rstp_interface.all",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_rstp_interface.all2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_rstp_interface.testacc_rstp_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRstpInterfaceConfigCreate(),
				},
				{
					Config: testAccJunosRstpInterfaceConfigUpdate(),
				},
				{
					ResourceName:      "junos_rstp_interface.all",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosRstpInterfaceSWConfigCreate() string {
	return `
resource "junos_rstp_interface" "all" {
  name = "all"
}

resource "junos_routing_instance" "testacc_rstp_interface" {
  name = "testacc_rstp_intface"
  type = "virtual-switch"
}

resource "junos_rstp_interface" "all2" {
  name             = "all"
  routing_instance = junos_routing_instance.testacc_rstp_interface.name
}
`
}

func testAccJunosRstpInterfaceSWConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_rstp_interface" "all" {
  name                      = "all"
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}

resource "junos_interface_physical" "testacc_rstp_interface" {
  name         = "%s"
  vlan_members = ["default"]
}

resource "junos_rstp_interface" "testacc_rstp_interface" {
  name         = junos_interface_physical.testacc_rstp_interface.name
  no_root_port = true
}

resource "junos_routing_instance" "testacc_rstp_interface" {
  name = "testacc_rstp_interface"
  type = "virtual-switch"
}

resource "junos_rstp_interface" "all2" {
  name                      = "all"
  routing_instance          = junos_routing_instance.testacc_rstp_interface.name
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}
`, interFace)
}

func testAccJunosRstpInterfaceConfigCreate() string {
	return `
resource "junos_rstp_interface" "all" {
  name = "all"
}
`
}

func testAccJunosRstpInterfaceConfigUpdate() string {
	return `
resource "junos_rstp_interface" "all" {
  name                      = "all"
  access_trunk              = true
  bpdu_timeout_action_alarm = true
  bpdu_timeout_action_block = true
  cost                      = 16
  edge                      = true
  mode                      = "shared"
  priority                  = 240
}
`
}
