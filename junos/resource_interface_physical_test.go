package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE_AE=ae<num> for choose interface aggregate test else it's ae0.
func TestAccJunosInterfacePhysical_basic(t *testing.T) {
	var testaccInterface string
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
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfacePhysicalConfigCreate(testaccInterface),
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
					Config: testAccJunosInterfacePhysicalConfigUpdate(testaccInterface),
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
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfacePhysicalFWConfigCreate(testaccInterface, testaccInterfaceAE),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interface"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"ether802_3ad", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"name", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"ae_lacp", "active"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"ae_minimum_links", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"vlan_tagging", "true"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.0.identifier", "00:01:11:11:11:11:11:11:11:11"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.0.mode", "all-active"),
					),
				},
				{
					Config: testAccJunosInterfacePhysicalFWConfigUpdate(testaccInterface, testaccInterfaceAE),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interfaceU"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"ae_lacp", ""),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"ae_minimum_links", "0"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.0.identifier", "00:11:11:11:11:11:11:11:11:11"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.0.mode", "all-active"),
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
			},
		})
	}
}

func testAccJunosInterfacePhysicalConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interface"
  trunk        = true
  vlan_native  = 100
  vlan_members = ["100-110"]
}
`)
}
func testAccJunosInterfacePhysicalConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interfaceU"
  vlan_members = ["100"]
}
`)
}

func testAccJunosInterfacePhysicalFWConfigCreate(interFace, interfaceAE string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interface"
  ether802_3ad = "` + interfaceAE + `"
}
resource junos_interface_physical testacc_interfaceAE {
  name             = junos_interface_physical.testacc_interface.ether802_3ad
  description      = "testacc_interfaceAE"
  ae_lacp          = "active"
  ae_minimum_links = 1
  vlan_tagging     = true
  esi {
    identifier = "00:01:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
}
`)
}
func testAccJunosInterfacePhysicalFWConfigUpdate(interFace, interfaceAE string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interfaceU"
  ether802_3ad = "` + interfaceAE + `"
}
resource junos_interface_physical testacc_interfaceAE {
  name         = junos_interface_physical.testacc_interface.ether802_3ad
  description  = "testacc_interfaceAE"
  vlan_tagging = true
  esi {
    identifier = "00:11:11:11:11:11:11:11:11:11"
    mode       = "all-active"
  }
}
`)
}
