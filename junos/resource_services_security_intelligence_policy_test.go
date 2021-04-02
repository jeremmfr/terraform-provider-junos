package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServicesSecurityIntellPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesSecurityIntellPolicyConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_policy.testacc_svcSecIntelPolicy",
							"category.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesSecurityIntellPolicyConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_policy.testacc_svcSecIntelPolicy",
							"category.#", "2"),
					),
				},
				{
					ResourceName:      "junos_services_security_intelligence_policy.testacc_svcSecIntelPolicy",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosServicesSecurityIntellPolicyConfigCreate() string {
	return `
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelPolicy_CC" {
  name     = "testacc_svcSecIntelPolicy_CC"
  category = "CC"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_security_intelligence_policy" "testacc_svcSecIntelPolicy" {
  name = "testacc_svcSecIntelPolicy#1"
  category {
    name         = "CC"
    profile_name = junos_services_security_intelligence_profile.testacc_svcSecIntelPolicy_CC.name
  }
}
`
}
func testAccJunosServicesSecurityIntellPolicyConfigUpdate() string {
	return `
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelPolicy_CC" {
  name     = "testacc_svcSecIntelPolicy_CC"
  category = "CC"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelPolicy_IPFilter" {
  name     = "testacc_svcSecIntelPolicy_IPFilter"
  category = "IPFilter"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services_security_intelligence_policy" "testacc_svcSecIntelPolicy" {
  name = "testacc_svcSecIntelPolicy#1"
  category {
    name         = "CC"
    profile_name = junos_services_security_intelligence_profile.testacc_svcSecIntelPolicy_CC.name
  }
  category {
    name         = "IPFilter"
    profile_name = junos_services_security_intelligence_profile.testacc_svcSecIntelPolicy_IPFilter.name
  }
}
`
}
