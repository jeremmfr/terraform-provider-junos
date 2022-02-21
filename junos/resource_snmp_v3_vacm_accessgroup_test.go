package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSnmpV3VacmAccessGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSnmpV3VacmAccessGroupConfigCreate(),
			},
			{
				Config: testAccJunosSnmpV3VacmAccessGroupConfigUpdate(),
			},
			{
				ResourceName:      "junos_snmp_v3_vacm_accessgroup.testacc_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSnmpV3VacmAccessGroupConfigCreate() string {
	return `
resource "junos_snmp_v3_vacm_accessgroup" "testacc_group" {
  name = "testacc_group#1"
  default_context_prefix {
    model       = "any"
    level       = "none"
    notify_view = "all"
  }
  context_prefix {
    prefix = "ctx#22"
    access_config {
      model     = "any"
      level     = "authentication"
      read_view = "all"
    }
  }
}
resource "junos_snmp_v3_vacm_accessgroup" "testacc_group2" {
  name = "testacc_group#2"
  context_prefix {
    prefix = "ctx"
    access_config {
      model       = "any"
      level       = "none"
      notify_view = "all"
    }
  }
}
`
}

func testAccJunosSnmpV3VacmAccessGroupConfigUpdate() string {
	return `
resource "junos_snmp_v3_vacm_accessgroup" "testacc_group" {
  name = "testacc_group#1"
  default_context_prefix {
    model     = "any"
    level     = "authentication"
    read_view = "all"
  }
  default_context_prefix {
    model       = "any"
    level       = "none"
    notify_view = "all"
  }
  default_context_prefix {
    model         = "usm"
    level         = "privacy"
    context_match = "exact"
    notify_view   = "all"
    read_view     = "all"
    write_view    = "all"
  }
  context_prefix {
    prefix = "ctx#22"
    access_config {
      model     = "any"
      level     = "authentication"
      read_view = "all"
    }
    access_config {
      model       = "any"
      level       = "none"
      notify_view = "all"
    }
    access_config {
      model         = "usm"
      level         = "privacy"
      context_match = "exact"
      notify_view   = "all"
      read_view     = "all"
      write_view    = "al1"
    }
  }
  context_prefix {
    prefix = "ctx#21"
    access_config {
      context_match = "prefix"
      model         = "any"
      level         = "none"
      notify_view   = "all"
    }
  }
}
`
}
