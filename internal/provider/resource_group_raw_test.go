package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceGroupRaw_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				// 1
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				// 2
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				// 3
				ResourceName:      "junos_group_raw.testacc_foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// 4
				ResourceName:      "junos_group_raw.testacc_foo",
				ImportStateId:     "testacc foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// 5
				ResourceName:      "junos_group_raw.testacc_bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// 6
				ConfigDirectory:    config.TestStepDirectory(),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
