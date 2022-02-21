package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSnmpV3VacmSecurityToGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSnmpV3VacmSecurityToGroupConfigCreate(),
			},
			{
				Config: testAccJunosSnmpV3VacmSecurityToGroupConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_v3_vacm_securitytogroup.testacc_snmpv3secutogrp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSnmpV3VacmSecurityToGroupConfigCreate() string {
	return `
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp" {
  name  = "testacc_snmpv3secutogrp"
  model = "usm"
  group = "testacc_snmpv3secutogrp"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp2" {
  name  = "testacc_snmpv3secutogrp"
  model = "v1"
  group = "testacc_snmpv3secutogrp"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp3" {
  name  = "testacc_snmpv3secutogrp"
  model = "v2c"
  group = "testacc_snmpv3secutogrp"
}
`
}

func testAccJunosSnmpV3VacmSecurityToGroupConfigUpdate() string {
	return `
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp" {
  name  = "testacc_snmpv3secutogrp"
  model = "usm"
  group = "testacc_snmpv3secutogrp2"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp2" {
  name  = "testacc_snmpv3secutogrp"
  model = "v1"
  group = "testacc_snmpv3secutogrp2"
}
resource "junos_snmp_v3_vacm_securitytogroup" "testacc_snmpv3secutogrp3" {
  name  = "testacc_snmpv3secutogrp"
  model = "v2c"
  group = "testacc_snmpv3secutogrp2"
}
`
}
