package junos_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccDataSourceInterfacesPhysicalPresent_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = defaultInterfaceSwitchTestAcc
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfacesPhysicalPresentPreSwitch(testaccInterface),
				},
				{
					Config:  testAccDataSourceInterfacesPhysicalPresentPreSwitch(testaccInterface),
					Destroy: true,
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfig(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", testaccInterface),
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", "dsc"),
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", "lo0"),
					),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth",
							"interface_names.*", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.0.name", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.0.admin_status", "up"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003AdmUp",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentLo0",
							"interface_names.#", "1"),
					),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch2(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch2(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003AdmUp",
							"interface_names.#", "0"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfacesPhysicalPresentPreSwitch(testaccInterface),
				},
				{
					Config:  testAccDataSourceInterfacesPhysicalPresentPreSwitch(testaccInterface),
					Destroy: true,
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfig(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", testaccInterface),
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", "dsc"),
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresent",
							"interface_names.*", "lo0"),
					),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth",
							"interface_names.*", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.0.name", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.0.admin_status", "up"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_statuses.0.oper_status", "down"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003AdmUp",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003OperUp",
							"interface_names.#", "0"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentLo0",
							"interface_names.#", "1"),
					),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch2(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfacesPhysicalPresentConfigMatch2(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003AdmUp",
							"interface_names.#", "0"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceInterfacesPhysicalPresentPreSwitch(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_dataIfacesPhysPresent {
  name        = "%s"
  description = "testacc_dataIfacesPhysPresent"
}
resource junos_interface_logical testacc_dataIfacesPhysPresent {
  name        = "${junos_interface_physical.testacc_dataIfacesPhysPresent.name}.0"
  description = "testacc_dataIfacesPhysPresent"
}
`, interFace)
}

func testAccDataSourceInterfacesPhysicalPresentConfig(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_dataIfacesPhysPresent {
  name        = "%s"
  description = "testacc_dataIfacesPhysPresent"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresent {
}
`, interFace)
}

func testAccDataSourceInterfacesPhysicalPresentConfigMatch(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface_physical testacc_dataIfacesPhysPresent {
  name        = "%s"
  description = "testacc_dataIfacesPhysPresent"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentEth {
  match_name = "^%s-.*$"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentEth003 {
  match_name = "^%s$"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentEth003AdmUp {
  match_name     = "^%s$"
  match_admin_up = true
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentEth003OperUp {
  match_name    = "^%s$"
  match_oper_up = true
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentLo0 {
  match_name    = "^lo0$"
  match_oper_up = true
}
`, interFace, strings.Split(interFace, `-`)[0], interFace, interFace, interFace)
}

func testAccDataSourceInterfacesPhysicalPresentConfigMatch2(interFace string) string {
	return fmt.Sprintf(`
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentEth003AdmUp {
  match_name     = "^%s$"
  match_admin_up = true
}
`, interFace)
}
