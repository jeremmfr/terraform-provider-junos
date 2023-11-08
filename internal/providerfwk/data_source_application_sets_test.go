package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceApplicationSets_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
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
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_ssh_without_telnet",
							"application_sets.#", "0"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_ssh_with_telnet",
							"application_sets.#", "1"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_name",
							"application_sets.#", "2"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_appsets",
							"application_sets.#", "1"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_appsets",
							"application_sets.0.description", "test-data-source-appSet"),
					),
				},
			},
		})
	}
}
