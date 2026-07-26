package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceSnmpV3UsmUser_basic(t *testing.T) {
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
				ResourceName:      "junos_snmp_v3_usm_user.testacc_snmpv3user_2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "junos_snmp_v3_usm_user.testacc_snmpv3user_4",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory:    config.TestStepDirectory(),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"junos_snmp_v3_usm_user.testacc_snmpv3user_3",
							plancheck.ResourceActionUpdate,
						),
						plancheck.ExpectResourceAction(
							"junos_snmp_v3_usm_user.testacc_snmpv3user_3_copy",
							plancheck.ResourceActionNoop,
						),
					},
				},
			},
			{
				ConfigDirectory:    config.TestStepDirectory(),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(
							"junos_snmp_v3_usm_user.testacc_snmpv3user_3",
							plancheck.ResourceActionNoop,
						),
						plancheck.ExpectResourceAction(
							"junos_snmp_v3_usm_user.testacc_snmpv3user_3_copy",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
		},
	})
}

func TestAccResourceSnmpV3UsmUser_writeOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"authentication_key"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"authentication_key_wo"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"authentication_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"privacy_key"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"privacy_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"authentication_password"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"authentication_key"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"authentication_password_wo_version", "1"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"privacy_password"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"privacy_key"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"privacy_password_wo_version", "1"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"authentication_key"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"authentication_key_wo_version", "2"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo",
						"privacy_key_wo_version", "2"),
					resource.TestCheckNoResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"authentication_password"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"authentication_password_wo_version", "2"),
					resource.TestCheckResourceAttr("junos_snmp_v3_usm_user.testacc_snmpv3user_wo_2",
						"privacy_password_wo_version", "2"),
				),
			},
		},
	})
}
