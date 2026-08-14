package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSecurityAuthenticationKeyChain_basic(t *testing.T) {
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
					ResourceName:      "junos_security_authentication_key_chain.testacc_secauthKeyChain",
					ImportState:       true,
					ImportStateVerify: true,
					// on import, the date is read from the device without the leading zeros
					// removed by it, so it cannot match the value in the configuration
					ImportStateVerifyIgnore: []string{
						"key.0.start_time",
						"key.1.start_time",
					},
				},
				{
					ResourceName:      "junos_security_authentication_key_chain.testacc_secauthKeyChainAO",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
