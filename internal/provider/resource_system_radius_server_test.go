package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceSystemRadiusServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"address", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"secret", "password"),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ResourceName:             "junos_system_radius_server.testacc_radiusServer",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			// testing no_decode_secrets provider attribute
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
