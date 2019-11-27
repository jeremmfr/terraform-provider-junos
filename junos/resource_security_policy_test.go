package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosSecurityPolicy_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityPolicyConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_source_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_source_address.0", "testacc_address1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_destination_address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_destination_address.0", "any"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_application.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.match_application.0", "junos-ssh"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.log_init", "true"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.log_close", "true"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.0.count", "true"),
					),
				},
				{
					Config: testAccJunosSecurityPolicyConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.#", "2"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.1.then", "reject"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_securityPolicy",
							"policy.1.match_source_address.0", "testacc_address1"),
					),
				},
				{
					ResourceName:      "junos_security_policy.testacc_securityPolicy",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityPolicyConfigCreate() string {
	return fmt.Sprintf(`
resource junos_security_policy testacc_securityPolicy {
  from_zone = junos_security_zone.testacc_seczonePolicy1.name
  to_zone = junos_security_zone.testacc_seczonePolicy1.name
  policy {
    name = "testacc_Policy_1"
    match_source_address = [ "testacc_address1" ]
    match_destination_address = [ "any" ]
    match_application = [ "junos-ssh" ]
    log_init = true
    log_close = true
    count = true
  }
}

resource junos_security_zone testacc_seczonePolicy1 {
	name = "testacc_seczonePolicy1"
	address_book {
         name = "testacc_address1"
         network = "192.0.2.0/25"
       }
}
`)
}
func testAccJunosSecurityPolicyConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_security_policy testacc_securityPolicy {
  from_zone = junos_security_zone.testacc_seczonePolicy1.name
  to_zone = junos_security_zone.testacc_seczonePolicy1.name
  policy {
    name = "testacc_Policy_1"
    match_source_address = [ "testacc_address1" ]
    match_destination_address = [ "any" ]
    match_application = [ "junos-ssh" ]
    log_init = true
    log_close = true
    count = true
  }
  policy {
    name = "testacc_Policy_2"
    match_source_address = [ "testacc_address1" ]
    match_destination_address = [ "any" ]
    match_application = [ "any" ]
    then = "reject"
  }
}

resource junos_security_zone testacc_seczonePolicy1 {
	name = "testacc_seczonePolicy1"
	address_book {
         name = "testacc_address1"
         network = "192.0.2.0/25"
       }
}
`)
}
