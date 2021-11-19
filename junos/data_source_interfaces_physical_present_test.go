package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfacesPhysicalPresent_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
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
					resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE",
						"interface_names.*", testaccInterface),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003",
						"interface_names.#", "1"),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003",
						"interface_statuses.#", "1"),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003",
						"interface_statuses.0.name", testaccInterface),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003",
						"interface_statuses.0.admin_status", "up"),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003",
						"interface_statuses.0.oper_status", "down"),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003AdmUp",
						"interface_names.#", "1"),
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003OperUp",
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
					resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentGE003AdmUp",
						"interface_names.#", "0"),
				),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}

func testAccDataSourceInterfacesPhysicalPresentConfig(interFace string) string {
	return `
resource junos_interface_physical testacc_dataIfacesPhysPresent {
  name        = "` + interFace + `"
  description = "testacc_dataIfacesPhysPresent"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresent {
}
`
}

func testAccDataSourceInterfacesPhysicalPresentConfigMatch(interFace string) string {
	return `
resource junos_interface_physical testacc_dataIfacesPhysPresent {
  name        = "` + interFace + `"
  description = "testacc_dataIfacesPhysPresent"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentGE {
  match_name = "^ge-.*$"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentGE003 {
  match_name = "^` + interFace + `$"
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentGE003AdmUp {
  match_name     = "^` + interFace + `$"
  match_admin_up = true
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentGE003OperUp {
  match_name    = "^` + interFace + `$"
  match_oper_up = true
}
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentLo0 {
  match_name    = "^lo0$"
  match_oper_up = true
}
`
}

func testAccDataSourceInterfacesPhysicalPresentConfigMatch2(interFace string) string {
	return `
data junos_interfaces_physical_present testacc_dataIfacesPhysPresentGE003AdmUp {
  match_name     = "^` + interFace + `$"
  match_admin_up = true
}
`
}
