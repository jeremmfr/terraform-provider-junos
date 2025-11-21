package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccActionLoadConfig_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck: func() { testAccPreCheck(t) },
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				// tfversion.SkipBelow(tfversion.Version1_14_0),
				tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0"))),
			},
			Steps: []resource.TestStep{
				{
					// 1
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
				},
				{
					// 2
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_applications.testacc",
							"applications.#", "1"),
					),
				},
				{
					// 3
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("data.junos_applications.testacc",
							"applications.0.source_port"),
					),
				},
				{
					// 4
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					ConfigDirectory:          config.TestStepDirectory(),
				},
				{
					// 5
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
