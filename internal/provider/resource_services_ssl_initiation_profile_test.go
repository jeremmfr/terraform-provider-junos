package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServicesSSLInitiationProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_ssl_initiation_profile.testacc_sslInitProf",
							"custom_ciphers.#", "1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_ssl_initiation_profile.testacc_sslInitProf",
							"custom_ciphers.#", "2"),
					),
				},
				{
					ResourceName:      "junos_services_ssl_initiation_profile.testacc_sslInitProf",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
