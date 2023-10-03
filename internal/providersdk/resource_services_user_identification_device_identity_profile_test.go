package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceServicesUserIdentDeviceIdentityProfile_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceServicesUserIdentDeviceIdentityProfileCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_user_identification_device_identity_profile.testacc_devidProf",
							"attribute.#", "1"),
					),
				},
				{
					Config: testAccResourceServicesUserIdentDeviceIdentityProfileUpdate(),
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

func testAccResourceServicesUserIdentDeviceIdentityProfileCreate() string {
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

func testAccResourceServicesUserIdentDeviceIdentityProfileUpdate() string {
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
