package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityIdpCustomAttackGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityIdpCustomAttackGroupConfigCreate(),
				},
				{
					Config: testAccJunosSecurityIdpCustomAttackGroupConfigUpdate(),
				},
				{
					ResourceName:      "junos_security_idp_custom_attack_group.testacc_idpCustomAttackGroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityIdpCustomAttackGroupConfigCreate() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttackGroup" {
  name               = "testacc/#1_Group"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
resource "junos_security_idp_custom_attack_group" "testacc_idpCustomAttackGroup" {
  name = "testacc/#1_CustomAttackGroup"
  member = [
    junos_security_idp_custom_attack.testacc_idpCustomAttackGroup.name,
  ]
}
`
}

func testAccJunosSecurityIdpCustomAttackGroupConfigUpdate() string {
	return `
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttackGroup" {
  name               = "testacc/#1_Group"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
resource "junos_security_idp_custom_attack" "testacc_idpCustomAttackGroup2" {
  name               = "testacc/#1_Group2"
  recommended_action = "ignore"
  severity           = "info"
  attack_type_anomaly {
    direction = "any"
    service   = "TELNET"
    test      = "SUBOPTION_OVERFLOW"
    shellcode = "all"
  }
}
resource "junos_security_idp_custom_attack_group" "testacc_idpCustomAttackGroup" {
  name = "testacc/#1_CustomAttackGroup"
  member = [
    junos_security_idp_custom_attack.testacc_idpCustomAttackGroup.name,
    junos_security_idp_custom_attack.testacc_idpCustomAttackGroup2.name,
  ]
}
`
}
