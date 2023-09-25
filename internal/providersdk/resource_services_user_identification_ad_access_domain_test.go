package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceServicesUserIdentAdAccessDomain_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceServicesUserIdentAdAccessDomainCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"domain_controller.#", "1"),
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"ip_user_mapping_discovery_wmi.#", "1"),
						resource.TestCheckResourceAttr("junos_services_user_identification_ad_access_domain.testacc_userID_addomain",
							"user_group_mapping_ldap.#", "1"),
					),
				},
				{
					Config: testAccResourceServicesUserIdentAdAccessDomainUpdate(),
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
				{
					Config: testAccResourceServicesUserIdentAdAccessDomainPostCheck(),
				},
			},
		})
	}
}

func testAccResourceServicesUserIdentAdAccessDomainCreate() string {
	return `
resource "junos_services" "testacc_userID_addomain" {
  user_identification {
    ad_access {}
  }
}
resource "junos_services_user_identification_ad_access_domain" "testacc_userID_addomain" {
  name          = "testacc_userID_addomain.local"
  user_name     = "user_dom"
  user_password = "user_pass"
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
  ip_user_mapping_discovery_wmi {}
  user_group_mapping_ldap {
    base = "CN=xxx"
  }
}
`
}

func testAccResourceServicesUserIdentAdAccessDomainUpdate() string {
	return `
resource "junos_services" "testacc_userID_addomain" {
  user_identification {
    ad_access {}
  }
}
resource "junos_services_user_identification_ad_access_domain" "testacc_userID_addomain" {
  name          = "testacc_userID_addomain.local"
  user_name     = "user_dom"
  user_password = "user_pass"
  domain_controller {
    name    = "server1"
    address = "192.0.2.3"
  }
  domain_controller {
    name    = "server0"
    address = "192.0.2.2"
  }
  ip_user_mapping_discovery_wmi {
    event_log_scanning_interval = 30
    initial_event_log_timespan  = 30
  }
  user_group_mapping_ldap {
    base             = "CN=xxx"
    address          = ["192.0.2.6", "192.0.2.5"]
    auth_algo_simple = true
    ssl              = true
    user_name        = "user_ldap_map"
    user_password    = "user_ldap_pass"
  }
}
`
}

func testAccResourceServicesUserIdentAdAccessDomainPostCheck() string {
	return `
resource "junos_services" "testacc_userID_addomain" {
}
`
}
