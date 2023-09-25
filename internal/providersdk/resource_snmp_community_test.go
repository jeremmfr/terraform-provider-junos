package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSnmpCommunity_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSnmpCommunityConfigCreate(),
			},
			{
				Config: testAccResourceSnmpCommunityConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_community.testacc_snmpcom",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSnmpCommunityConfigCreate() string {
	return `
resource "junos_snmp" "testacc_snmpcom" {
  routing_instance_access = true
}
resource "junos_snmp_clientlist" "testacc_snmpcom" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_snmpcom"
}
resource "junos_snmp_community" "testacc_snmpcom" {
  depends_on = [
    junos_snmp.testacc_snmpcom
  ]
  name                    = "testacc_snmpcom@public"
  authorization_read_only = true
  client_list_name        = junos_snmp_clientlist.testacc_snmpcom.name
  routing_instance {
    name = junos_routing_instance.testacc_snmpcom.name
  }
  view = junos_snmp_view.testacc_snmpcom.name
}
resource "junos_routing_instance" "testacc_snmpcom" {
  name = "testacc_snmpcom"
}
resource "junos_snmp_view" "testacc_snmpcom" {
  lifecycle {
    create_before_destroy = true
  }
  name        = "testacc_snmpcom"
  oid_include = [".1"]
}
`
}

func testAccResourceSnmpCommunityConfigUpdate() string {
	return `
resource "junos_snmp_community" "testacc_snmpcom" {
  name                     = "testacc_snmpcom@public"
  authorization_read_write = true
  clients                  = ["192.0.2.0/24"]
  routing_instance {
    name = junos_routing_instance.testacc_snmpcom.name
  }
  routing_instance {
    name = junos_routing_instance.testacc2_snmpcom.name
  }
}

resource "junos_routing_instance" "testacc_snmpcom" {
  name = "testacc_snmpcom"
}
resource "junos_routing_instance" "testacc2_snmpcom" {
  name = "testacc2_snmpcom"
}
`
}
