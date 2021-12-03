package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServicesSSLInitiationProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesSSLInitiationProfileCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_ssl_initiation_profile.testacc_sslInitProf",
							"actions.#", "1"),
						resource.TestCheckResourceAttr("junos_services_ssl_initiation_profile.testacc_sslInitProf",
							"custom_ciphers.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesSSLInitiationProfileUpdate(),
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

func testAccJunosServicesSSLInitiationProfileCreate() string {
	return `
resource "junos_services_ssl_initiation_profile" "testacc_sslInitProf" {
  name = "testacc_sslInitProf.1"
  actions {
    crl_disable                      = true
    crl_if_not_present               = "allow"
    crl_ignore_hold_instruction_code = true
    ignore_server_auth_failure       = true
  }
  custom_ciphers       = ["rsa-with-aes-128-gcm-sha256"]
  enable_flow_tracing  = true
  enable_session_cache = true
  preferred_ciphers    = "medium"
  protocol_version     = "all"
  # trusted_ca         = ["all"] # fail on recent version of Junos
}
`
}

func testAccJunosServicesSSLInitiationProfileUpdate() string {
	return `
resource "junos_services_ssl_initiation_profile" "testacc_sslInitProf" {
  name = "testacc_sslInitProf.1"
  actions {
    crl_disable                      = true
    crl_if_not_present               = "allow"
    crl_ignore_hold_instruction_code = true
    ignore_server_auth_failure       = true
  }
  custom_ciphers       = ["rsa-with-aes-128-gcm-sha256", "rsa-with-aes-256-cbc-sha"]
  enable_flow_tracing  = true
  enable_session_cache = true
  preferred_ciphers    = "medium"
  protocol_version     = "tls12"
  # trusted_ca         = ["all"] # fail on recent version of Junos
}
`
}
