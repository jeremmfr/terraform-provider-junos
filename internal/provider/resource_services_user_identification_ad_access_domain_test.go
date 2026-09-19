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

func TestAccResourceServicesUserIdentificationADAccessDomain_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"domain_controller.#", "1"),
						resource.TestCheckResourceAttrSet("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"ip_user_mapping_discovery_wmi.%"),
						resource.TestCheckResourceAttrSet("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"user_group_mapping_ldap.base"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"domain_controller.#", "2"),
					),
				},
				{
					ResourceName:      "junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func TestAccResourceServicesUserIdentificationADAccessDomain_writeOnly(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					// 1
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_password"),
						resource.TestCheckNoResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_password_wo"),
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_password_wo_version", "1"),
						resource.TestCheckNoResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_group_mapping_ldap.user_password"),
						resource.TestCheckNoResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_group_mapping_ldap.user_password_wo"),
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_group_mapping_ldap.user_password_wo_version", "1"),
					),
				},
				{
					// 2 check that the write-only passwords have really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_userID_addomain_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`domain "?testacc_userID_addomain_wo\.local"? user password `,
							)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_userID_addomain_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`user-group-mapping ldap user password `,
							)),
						),
					},
				},
				{
					// 3
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_password_wo_version", "2"),
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain_wo",
							"user_group_mapping_ldap.user_password_wo_version", "2"),
					),
				},
			},
		})
	}
}
