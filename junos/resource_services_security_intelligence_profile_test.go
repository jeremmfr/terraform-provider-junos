package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServicesSecurityIntellProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesSecurityIntellProfileConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
							"rule.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesSecurityIntellProfileConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
							"rule.#", "4"),
					),
				},
				{
					ResourceName:      "junos_services_security_intelligence_profile.testacc_svcSecIntelProfile",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosServicesSecurityIntellProfileConfigCreate() string {
	return `
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelProfile" {
  name     = "testacc_svcSecIntelProfile@1"
  category = "CC"
  rule {
    name = "test#2"
    match {
      threat_level = [10]
      feed_name    = ["CC_IP"]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
}
`
}

func testAccJunosServicesSecurityIntellProfileConfigUpdate() string {
	return `
resource "junos_services_security_intelligence_profile" "testacc_svcSecIntelProfile" {
  name     = "testacc_svcSecIntelProfile@1"
  category = "CC"
  default_rule_then {
    action = "permit"
    no_log = true
  }
  rule {
    name = "test#3"
    match {
      threat_level = [5, 4]
      feed_name    = ["CC_URL"]
    }
    then_action = "permit"
    then_log    = true
  }
  rule {
    name = "test"
    match {
      threat_level = [1]
    }
    then_action = "recommended"
  }
  rule {
    name = "test#2"
    match {
      threat_level = [10]
      feed_name    = ["CC_IP"]
    }
    then_action = "block close http redirect-url http://www.test.com/url1.html"
    then_log    = true
  }
  rule {
    name = "test2"
    match {
      threat_level = [10]
    }
    then_action = "block drop"
  }
}
`
}
