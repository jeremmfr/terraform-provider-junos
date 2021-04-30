package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServicesUserIdentDeviceIdentityProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesUserIdentDeviceIdentityProfileCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_device_identity_profile.testacc_devidProf",
							"attribute.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesUserIdentDeviceIdentityProfileUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_device_identity_profile.testacc_devidProf",
							"attribute.#", "2"),
					),
				},
				{
					ResourceName:      "junos_services_user_identification_device_identity_profile.testacc_devidProf",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosServicesUserIdentDeviceIdentityProfileCreate() string {
	return `
resource "junos_services_user_identification_device_identity_profile" "testacc_devidProf" {
  name   = "testacc_devidProf.1"
  domain = "clearpass"
  attribute {
    name  = "device-identity"
    value = ["device1", "barcode scan"]
  }
}
`
}

func testAccJunosServicesUserIdentDeviceIdentityProfileUpdate() string {
	return `
resource "junos_services_user_identification_device_identity_profile" "testacc_devidProf" {
  name   = "testacc_devidProf.1"
  domain = "clearpass"
  attribute {
    name  = "device-identity"
    value = ["device1", "barcode scan"]
  }
  attribute {
    name  = "device-category"
    value = ["category@1"]
  }
}
`
}
