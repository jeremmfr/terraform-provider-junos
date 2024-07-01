package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRipGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
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
					ResourceName:      "junos_rip_group.testacc_ripgroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_rip_group.testacc_ripgroup2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_rip_group.testacc_ripnggroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_rip_group.testacc_ripnggroup2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
