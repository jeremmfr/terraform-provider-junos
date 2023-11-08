package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceApplications_basic(t *testing.T) {
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
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_any",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_any",
							"applications.0.name", "any"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh-name",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh-name",
							"applications.0.name", "junos-ssh"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh",
							"applications.0.name", "junos-ssh"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_all_ssh",
							"applications.#", "3"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_multi_term",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_multi_term",
							"applications.0.name", "testacc_custom_multi_term"),
					),
				},
			},
		})
	}
}
