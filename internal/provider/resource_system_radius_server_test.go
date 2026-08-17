package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceSystemRadiusServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"address", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"secret", "password"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"preauthentication_secret", "password"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"source_address", "192.0.2.2"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"port", "1645"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_port", "1646"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"dynamic_request_port", "3799"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"preauthentication_port", "1812"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"timeout", "10"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_timeout", "5"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"retry", "3"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_retry", "2"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"max_outstanding_requests", "1000"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"routing_instance", "testacc_radiusServer"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ResourceName:             "junos_system_radius_server.testacc_radiusServer",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			// testing no_decode_secrets provider attribute
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccResourceSystemRadiusServer_writeOnly(t *testing.T) {
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
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret"),
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret_wo"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret_wo_version", "1"),
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"preauthentication_secret"),
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"preauthentication_secret_wo"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"preauthentication_secret_wo_version", "1"),
				),
			},
			{
				// check that the write-only secrets have really been sent to the device
				ConfigDirectory: config.TestStepDirectory(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.testacc_radiusServer_wo",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile(
							`set system radius-server 192.0.2.11 secret `)),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.testacc_radiusServer_wo",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile(
							`set system radius-server 192.0.2.11 preauthentication-secret `)),
					),
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret"),
					resource.TestCheckNoResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret_wo"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"secret_wo_version", "2"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer_wo",
						"preauthentication_secret_wo_version", "2"),
				),
			},
		},
	})
}
