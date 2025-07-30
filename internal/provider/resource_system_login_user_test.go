package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccResourceSystemLoginUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_login_user.testacc",
						"name", "testacc"),
					resource.TestCheckResourceAttrSet("junos_system_login_user.testacc",
						"uid"),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"junos_system_login_user.testacc2",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_login_user.testacc",
						"name", "testacc"),
					resource.TestCheckResourceAttrSet("junos_system_login_user.testacc",
						"uid"),
					resource.TestCheckResourceAttr("junos_system_login_user.testacc",
						"authentication.ssh_public_keys.#", "1"),
					resource.TestCheckResourceAttr("junos_system_login_user.testacc2",
						"uid", "5000"),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ResourceName:             "junos_system_login_user.testacc",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				ExpectNonEmptyPlan:       true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"junos_system_login_user.testacc3",
							plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectResourceAction(
							"junos_system_login_user.testacc3_copy",
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
							"junos_system_login_user.testacc3",
							plancheck.ResourceActionNoop,
						),
						plancheck.ExpectResourceAction(
							"junos_system_login_user.testacc3_copy",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
		},
	})
}
