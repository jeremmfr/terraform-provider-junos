package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceApplyGroupExcept_basic(t *testing.T) {
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
				ResourceName:      "junos_apply_group_except.testacc_foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
