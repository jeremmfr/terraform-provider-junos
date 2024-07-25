package providerfwk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
