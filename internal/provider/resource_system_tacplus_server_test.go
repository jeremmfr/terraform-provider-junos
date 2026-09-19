package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceSystemTacplusServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ResourceName:      "junos_system_tacplus_server.testacc_tacplusServer",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
		},
	})
}

func TestAccResourceSystemTacplusServer_writeOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret"),
					resource.TestCheckNoResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret_wo"),
					resource.TestCheckResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret_wo_version", "1"),
				),
			},
			{
				// check that the write-only secret has really been sent to the device
				ConfigDirectory: config.TestStepDirectory(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.testacc_tacplusServer_wo",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile(
							`set system tacplus-server 192.0.2.11 secret `,
						)),
					),
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret"),
					resource.TestCheckNoResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret_wo"),
					resource.TestCheckResourceAttr("junos_system_tacplus_server.testacc_tacplusServer_wo",
						"secret_wo_version", "2"),
				),
			},
		},
	})
}
