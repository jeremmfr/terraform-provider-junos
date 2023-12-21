package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's ge-0/0/4.
// export TESTACC_INTERFACE_AE=ae<num> for choose interface aggregate test else it's ae0.
func TestAccResourceInterfacePhysical_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	testaccInterfaceAE := "ae0"
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE_AE"); iface != "" {
		testaccInterfaceAE = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interface2":  config.StringVariable(testaccInterface2),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interface"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"gigether_opts.ae_8023ad", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"name", testaccInterfaceAE),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.lacp.mode", "active"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.minimum_links", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"vlan_tagging", "true"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interface2":  config.StringVariable(testaccInterface2),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interfaceU"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.lacp.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.minimum_bandwidth", "1 gbps"),
					),
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interface2":  config.StringVariable(testaccInterface2),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					ResourceName:      "junos_interface_physical.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interface2":  config.StringVariable(testaccInterface2),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					ResourceName:      "junos_interface_physical.testacc_interfaceAE",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interface2":  config.StringVariable(testaccInterface2),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
				},
			},
		})
	}
}

func TestAccResourceInterfacePhysical_router(t *testing.T) {
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
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.source_address_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"parent_ether_opts.source_filtering", "true"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.identifier", "00:01:11:11:11:11:11:11:11:11"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.mode", "all-active"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.identifier", "00:11:11:11:11:11:11:11:11:11"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interfaceAE",
							"esi.mode", "all-active"),
					),
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface":   config.StringVariable(testaccInterface),
						"interfaceAE": config.StringVariable(testaccInterfaceAE),
					},
					ResourceName:      "junos_interface_physical.testacc_interfaceAE",
					ImportState:       true,
					ImportStateVerify: true,
				},
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
				},
			},
		})
	}
}

func TestAccResourceInterfacePhysical_srx(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface_reth",
							"parent_ether_opts.redundancy_group", "1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
			},
		})
	}
}

func TestAccResourceInterfacePhysical_switch(t *testing.T) {
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"description", "testacc_interfaceU"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_physical.testacc_interface",
							"vlan_members.0", "100"),
					),
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_interface_physical.testacc_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
