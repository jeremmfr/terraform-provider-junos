package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceSystemRootAuthentication_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() { testAccPreCheck(t) },
			Steps: []resource.TestStep{
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_root_authentication.root_auth",
							"encrypted_password", "$6$XXXX"),
						resource.TestCheckResourceAttr("junos_system_root_authentication.root_auth",
							"ssh_public_keys.#", "1"),
					),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ResourceName:             "junos_system_root_authentication.root_auth",
					ImportState:              true,
					ImportStateVerify:        true,
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_root_authentication.root_auth",
							"ssh_public_keys.#", "0"),
					),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					ExpectNonEmptyPlan:       true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							plancheck.ExpectNonEmptyPlan(),
							plancheck.ExpectResourceAction(
								"junos_system_root_authentication.root_auth",
								plancheck.ResourceActionUpdate,
							),
							plancheck.ExpectResourceAction(
								"junos_system_root_authentication.root_auth_copy",
								plancheck.ResourceActionNoop,
							),
						},
					},
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					ExpectNonEmptyPlan:       true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							plancheck.ExpectNonEmptyPlan(),
							plancheck.ExpectResourceAction(
								"junos_system_root_authentication.root_auth",
								plancheck.ResourceActionNoop,
							),
							plancheck.ExpectResourceAction(
								"junos_system_root_authentication.root_auth_copy",
								plancheck.ResourceActionUpdate,
							),
						},
					},
				},
			},
		})
	}
}
