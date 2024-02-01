package providerfwk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSnmpV3UsmUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ResourceName:      "junos_snmp_v3_usm_user.testacc_snmpv3user_2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "junos_snmp_v3_usm_user.testacc_snmpv3user_4",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
