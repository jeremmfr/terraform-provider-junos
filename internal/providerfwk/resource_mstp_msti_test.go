package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceMstpMsti_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = junos.DefaultInterfaceSwitchTestAcc
	}
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface2 = junos.DefaultInterfaceSwitchTestAcc2
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"interface":  config.StringVariable(testaccInterface),
					"interface2": config.StringVariable(testaccInterface2),
				},
			},
			{
				ConfigVariables: map[string]config.Variable{
					"interface":  config.StringVariable(testaccInterface),
					"interface2": config.StringVariable(testaccInterface2),
				},
				ResourceName:      "junos_mstp_msti.testacc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's xe-0/0/3.
func TestAccResourceMstpMsti_switch(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceSwitchTestAcc
	testaccInterface2 := junos.DefaultInterfaceSwitchTestAcc2
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface":  config.StringVariable(testaccInterface),
						"interface2": config.StringVariable(testaccInterface2),
					},
					ResourceName:      "junos_mstp_msti.testacc",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
