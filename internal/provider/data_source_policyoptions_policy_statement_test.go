package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsPolicyStatement_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.junos_policyoptions_policy_statement.testacc_dataPolicyStatement",
						"id", "testacc_dataPolicyStatement"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_policy_statement.testacc_dataPolicyStatement2",
						"id", "testacc_dataPolicyStatement2"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_policy_statement.testacc_dataPolicyStatement3",
						"id", "testacc_dataPolicyStatement3"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options policy-statement .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
