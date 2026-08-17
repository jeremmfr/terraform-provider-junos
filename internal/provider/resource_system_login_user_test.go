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

func TestAccResourceSystemLoginUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_login_user.testacc",
						"name", "testacc"),
					resource.TestCheckResourceAttrSet("junos_system_login_user.testacc",
						"uid"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ResourceName:             "junos_system_login_user.testacc",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

func TestAccResourceSystemLoginUser_writeOnly(t *testing.T) {
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
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password"),
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password_wo"),
					resource.TestCheckResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password_wo_version", "1"),
				),
			},
			{
				// check that the write-only password has really been sent to the device
				ConfigDirectory: config.TestStepDirectory(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.testacc_wo",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile(
							`set system login user testacc_wo authentication encrypted-password `)),
					),
				},
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password"),
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password_wo"),
					resource.TestCheckResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password_wo_version", "2"),
				),
			},
			{
				// switch to a write-only plain text password,
				// the encrypted password generated by the device must not be kept in the private state
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password"),
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.plain_text_password"),
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.plain_text_password_wo"),
					resource.TestCheckResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.plain_text_password_wo_version", "1"),
					resource.TestCheckNoResourceAttr("junos_system_login_user.testacc_wo",
						"authentication.encrypted_password_wo_version"),
				),
			},
			{
				// the plain text password is unchanged, so the plan must stay empty
				// even though the provider cannot compare it with the device
				ConfigDirectory: config.TestStepDirectory(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
