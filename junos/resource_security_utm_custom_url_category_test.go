package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityUtmCustomURLCategory_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityUtmCustomURLCategoryConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_custom_url_category.testacc_URLCategory",
							"value.#", "1"),
						resource.TestCheckResourceAttr("junos_security_utm_custom_url_category.testacc_URLCategory",
							"value.0", "testacc-custom-pattern1"),
					),
				},
				{
					Config: testAccJunosSecurityUtmCustomURLCategoryConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_utm_custom_url_category.testacc_URLCategory",
							"value.#", "2"),
						resource.TestCheckResourceAttr("junos_security_utm_custom_url_category.testacc_URLCategory",
							"value.0", "testacc-custom-pattern1"),
						resource.TestCheckResourceAttr("junos_security_utm_custom_url_category.testacc_URLCategory",
							"value.1", "testacc-custom-pattern2"),
					),
				},
				{
					ResourceName:      "junos_security_utm_custom_url_category.testacc_URLCategory",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityUtmCustomURLCategoryConfigCreate() string {
	return `
resource junos_security_utm_custom_url_pattern "testacc_URLCategory1" {
  name  = "testacc-custom-pattern1"
  value = ["*.google.com"]
}
resource junos_security_utm_custom_url_category "testacc_URLCategory" {
  name = "testacc_URLCategory"
  value = [
    junos_security_utm_custom_url_pattern.testacc_URLCategory1.name,
  ]
}
`
}
func testAccJunosSecurityUtmCustomURLCategoryConfigUpdate() string {
	return `
resource junos_security_utm_custom_url_pattern "testacc_URLCategory1" {
  name  = "testacc-custom-pattern1"
  value = ["*.google.com"]
}
resource junos_security_utm_custom_url_pattern "testacc_URLCategory2" {
  name  = "testacc-custom-pattern2"
  value = ["*.google.fr"]
}
resource junos_security_utm_custom_url_category "testacc_URLCategory" {
  name = "testacc_URLCategory"
  value = [
    junos_security_utm_custom_url_pattern.testacc_URLCategory1.name,
    junos_security_utm_custom_url_pattern.testacc_URLCategory2.name,
  ]
}
`
}
