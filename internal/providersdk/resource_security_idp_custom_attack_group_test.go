package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityIdpCustomAttackGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityIdpCustomAttackGroupConfigCreate(),
				},
				{
					Config: testAccResourceSecurityIdpCustomAttackGroupConfigUpdate(),
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

func testAccResourceSecurityIdpCustomAttackGroupConfigCreate() string {
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

func testAccResourceSecurityIdpCustomAttackGroupConfigUpdate() string {
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
