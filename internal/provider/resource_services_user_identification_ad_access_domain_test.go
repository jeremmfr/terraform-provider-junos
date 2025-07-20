package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServicesUserIdentificationADAccessDomain_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
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
