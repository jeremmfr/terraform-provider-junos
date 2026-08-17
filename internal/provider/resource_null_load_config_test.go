package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceNullLoadConfig_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() { testAccPreCheck(t) },
			Steps: []resource.TestStep{
				{
					// 1
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
				},
				{
					// 2
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_applications.testacc",
							"applications.#", "1"),
					),
				},
				{
					// 3
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("data.junos_applications.testacc",
							"applications.0.source_port"),
					),
				},
				{
					// 4
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
				},
				{
					// 5
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_applications.testacc",
							"applications.#", "0"),
						resource.TestCheckResourceAttr("data.junos_applications.test_acc",
							"applications.#", "1"),
					),
				},
				{
					// 6
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_applications.testacc",
							"applications.#", "0"),
						resource.TestCheckResourceAttr("data.junos_applications.test_acc",
							"applications.#", "0"),
					),
				},
			},
		})
	}
}
