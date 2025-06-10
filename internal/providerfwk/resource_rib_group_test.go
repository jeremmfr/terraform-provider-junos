package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRibGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_policy.#", "1"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_policy.0", "testacc ribGroup"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.#", "1"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.0", "testacc_ribGroup1.inet.0"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"export_rib", "testacc_ribGroup1.inet.0"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.#", "3"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.2", "testacc_ribGroup2.inet.0"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"export_rib", "testacc_ribGroup2.inet.0"),
					),
				},
				{
					ResourceName:      "junos_rib_group.testacc_ribGroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
