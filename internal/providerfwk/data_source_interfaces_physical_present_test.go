package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccDataSourceInterfacesPhysicalPresent_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Destroy: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth",
							"interface_names.*", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces.%", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".name", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".admin_status", "up"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".oper_status", "down"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".logical_interface_names.#", "1"),
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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

func TestAccDataSourceInterfacesPhysicalPresent_switch(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceSwitchTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Destroy: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth",
							"interface_names.*", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interface_names.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces.%", "1"),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".name", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interfaces_physical_present.testacc_dataIfacesPhysPresentEth003",
							"interfaces."+testaccInterface+".admin_status", "up"),
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
