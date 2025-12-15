package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceApplyGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
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
					ResourceName:      "junos_apply_group.testacc_foobar",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 4
					ResourceName:      "junos_apply_group.testacc_barfoo",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 5
					ResourceName:      "junos_apply_group.testacc_barfoo2",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
