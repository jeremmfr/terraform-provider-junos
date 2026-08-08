package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceServices_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("junos_services.testacc",
							"security_intelligence.authentication_token"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("junos_services.testacc",
							"security_intelligence.authentication_token"),
						resource.TestCheckResourceAttr("junos_services.testacc",
							"security_intelligence.default_policy.#", "1"),
					),
				},
				{
					ResourceName:      "junos_services.testacc",
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

func TestAccResourceServices_writeOnly(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter_wo"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter_wo_version", "1"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret_wo"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret_wo_version", "1"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret_wo"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret_wo_version", "1"),
					),
				},
				{
					// check that the write-only secrets have really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set services security-intelligence url-parameter `)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set services user-identification identity-management connection primary client-secret `)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set services user-identification identity-management connection secondary client-secret `)),
						),
					},
				},
				{
					// increment the version of the url parameter and of the primary client secret
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter_wo_version", "2"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret_wo_version", "2"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret_wo_version", "1"),
					),
				},
				{
					// switch the url parameter and the secondary client secret
					// back to the standard arguments
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter", "test_param3"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"security_intelligence.url_parameter_wo_version"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.primary_client_secret_wo_version", "2"),
						resource.TestCheckResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret", "mySecret3"),
						resource.TestCheckNoResourceAttr("junos_services.testacc_wo",
							"user_identification.identity_management.connection.secondary_client_secret_wo_version"),
					),
				},
				{
					// clean the services configuration on the device
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
