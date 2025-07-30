package provider_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceMstpInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ResourceName:      "junos_mstp_interface.all",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's xe-0/0/3.
func TestAccResourceMstpInterface_switch(t *testing.T) {
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
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_mstp_interface.all",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_mstp_interface.testacc_mstp_interface",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
