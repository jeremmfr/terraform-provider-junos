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
func TestAccResourceOspfArea_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
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
					// 1
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"area_id", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"version", "v2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.name", "all"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.disable", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.passive", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.metric", "100"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.retransmit_interval", "12"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.hello_interval", "11"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.dead_interval", "10"),
					),
				},
				{
					// 2
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"routing_instance", "testacc_ospfarea"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"version", "v3"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea2",
							"interface.1.name", testaccInterface2+".0"),
					),
				},
				{
					// 3
					ResourceName:      "junos_ospf_area.testacc_ospfarea",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					// 4
					ResourceName:      "junos_ospf_area.testacc_ospfarea2",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					// 5
					ResourceName:      "junos_ospf_area.testacc_ospfareav3ipv4",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					// 6
					ResourceName:      "junos_ospf_area.testacc_ospfarea2v3realm",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					// 7
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 8
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 9
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					// 10
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					// 11
					ResourceName:      "junos_ospf_area.testacc_ospfarea",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					// 12
					ResourceName:      "junos_ospf_area.testacc_ospfarea2",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					// 13
					ResourceName:      "junos_ospf_area.testacc_ospfarea3",
					ImportState:       true,
					ImportStateVerify: true,
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
			},
		})
	}
}
