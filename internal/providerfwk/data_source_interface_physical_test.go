package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfacePhysical_basic(t *testing.T) {
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
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_datainterfaceP",
							"id", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_datainterfaceP",
							"vlan_tagging", "true"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func TestAccDataSourceInterfacePhysical_router(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterfaceAE := "ae0"
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE_AE"); iface != "" {
		testaccInterfaceAE = iface
	}
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_interfaceAE",
							"id", testaccInterfaceAE),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func TestAccDataSourceInterfacePhysical_switch(t *testing.T) {
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
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_interface",
							"id", testaccInterface),
						resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_interface",
							"storm_control", "testacc interface"),
					),
				},
			},
		})
	}
}
