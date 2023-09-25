package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityDynamicAddressName_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityDynamicAddressNameConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name",
							"profile_feed_name", "feedtfacc"),
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name2",
							"profile_category.#", "1"),
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name2",
							"profile_category.0.property.#", "2"),
					),
				},
				{
					ResourceName:      "junos_security_dynamic_address_name.testacc_dyn_add_name",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_dynamic_address_name.testacc_dyn_add_name2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceSecurityDynamicAddressNameConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_name.testacc_dyn_add_name",
							"profile_category.#", "1"),
					),
				},
			},
		})
	}
}

func testAccResourceSecurityDynamicAddressNameConfigCreate() string {
	return `
resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_name" {
  lifecycle {
    create_before_destroy = true
  }

  name     = "tfacc_dynadd"
  hostname = "example.com"
  feed_name {
    name = "feedtfacc"
    path = "/srx/"
  }
}
resource "junos_security_dynamic_address_name" "testacc_dyn_add_name" {
  name              = "tfacc_dynadd"
  description       = "desc tfacc dynamic-address"
  profile_feed_name = junos_security_dynamic_address_feed_server.testacc_dyn_add_name.feed_name.0.name
}
resource "junos_security_dynamic_address_name" "testacc_dyn_add_name2" {
  name        = "tfacc_dynadd2"
  description = "desc tfacc dynamic-address2"
  profile_category {
    name = "IPFilter"
    property {
      name   = "others#1"
      string = ["test#2", "test#1"]
    }
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
`
}

func testAccResourceSecurityDynamicAddressNameConfigUpdate() string {
	return `
resource "junos_security_dynamic_address_name" "testacc_dyn_add_name" {
  name        = "tfacc_dynadd"
  description = "desc tfacc dynamic-address"
  profile_category {
    name = "GeoIP"
    feed = "cat_feed"
    property {
      name   = "country"
      string = ["AU", "CN"]
    }
  }
}
`
}
