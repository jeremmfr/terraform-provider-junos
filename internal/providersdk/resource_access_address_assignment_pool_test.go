package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAccessAddressAssignmentPool_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceAccessAddressAssignmentPoolCreate(),
				},
				{
					Config: testAccResourceAccessAddressAssignmentPoolUpdate(),
				},
				{
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP6_1",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP6_2",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceAccessAddressAssignmentPoolCreate() string {
	return `
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP4" {
  name = "testacc_accessAddAssP4"
  family {
    type    = "inet"
    network = "192.0.2.128/25"
  }
}

resource "junos_routing_instance" "testacc_accessAddAssP6" {
  name = "testacc_accessAddAssP6"
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP6_1" {
  name             = "testacc_accessAddAssP6_1"
  routing_instance = junos_routing_instance.testacc_accessAddAssP6.name
  family {
    type    = "inet6"
    network = "fe80:0:0:b::/64"
  }
}
`
}

func testAccResourceAccessAddressAssignmentPoolUpdate() string {
	return `
resource "junos_interface_logical" "testacc_accessAddAssP4" {
  name = "lo0.1"
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP4" {
  name = "testacc_accessAddAssP4"
  family {
    type    = "inet"
    network = "192.0.2.128/25"
    excluded_address = [
      "192.0.2.201",
      "192.0.2.200",
    ]
    excluded_range {
      name = "excl1"
      low  = "192.0.2.208"
      high = "192.0.2.210"
    }
    excluded_range {
      name = "excl2"
      low  = "192.0.2.219"
      high = "192.0.2.220"
    }
    host {
      name             = "host1"
      hardware_address = "aa:bb:cc:dd:ee:ff"
      ip_address       = "192.0.2.211"
    }
    inet_range {
      name = "range2"
      low  = "192.0.2.225"
      high = "192.0.2.230"
    }
    inet_range {
      name = "range1"
      low  = "192.0.2.200"
      high = "192.0.2.220"
    }
    xauth_attributes_primary_dns    = "192.0.2.53/32"
    xauth_attributes_primary_wins   = "192.0.2.54/32"
    xauth_attributes_secondary_dns  = "192.0.2.55/32"
    xauth_attributes_secondary_wins = "192.0.2.56/32"
    dhcp_attributes {
      boot_file          = "file.boot"
      boot_server        = "test.com"
      domain_name        = "test.com"
      grace_period       = 3600
      maximum_lease_time = 3600
      name_server        = ["192.0.2.2"]
      netbios_node_type  = "b-node"
      next_server        = "192.0.2.3"
      option = [
        "1 string a",
        "2 flag true",
      ]
      option_match_82_circuit_id {
        value = "bb"
        range = "cc"
      }
      option_match_82_remote_id {
        value = "dd"
        range = "ee"
      }
      propagate_ppp_settings      = [junos_interface_logical.testacc_accessAddAssP4.name, ]
      propagate_settings          = "ff"
      router                      = ["192.0.2.121", "192.0.2.120"]
      server_identifier           = "192.0.2.6"
      sip_server_inet_address     = ["192.0.2.62", "192.0.2.61"]
      sip_server_inet_domain_name = ["domain.name"]
      t1_percentage               = 50
      t2_percentage               = 70
      tftp_server                 = "192.0.2.7"
      wins_server                 = ["192.0.2.72", "192.0.2.71"]
    }
  }
  active_drain = true
  hold_down    = true
}

resource "junos_routing_instance" "testacc_accessAddAssP6" {
  name = "testacc_accessAddAssP6"
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP6_1" {
  name             = "testacc_accessAddAssP6_1"
  routing_instance = junos_routing_instance.testacc_accessAddAssP6.name
  family {
    type    = "inet6"
    network = "fe80:0:0:b::/64"
    excluded_address = [
      "fe80:0:0:b::bb",
      "fe80:0:0:b::aa",
    ]
    inet6_range {
      name = "range62"
      low  = "fe80:0:0:b:1::bbbb/80"
      high = "fe80:0:0:b:1::cccc/80"
    }
    inet6_range {
      name = "range6"
      low  = "fe80:0:0:b:2::bbbb/80"
      high = "fe80:0:0:b:2::aaaa/80"
    }
    inet6_range {
      name          = "range_pref"
      prefix_length = 100
    }
    dhcp_attributes {
      dns_server                   = ["fe80::1"]
      exclude_prefix_len           = "65"
      maximum_lease_time_infinite  = true
      sip_server_inet6_address     = ["fe80:0:0:b:b::cc", "fe80:0:0:b:b::dc"]
      sip_server_inet6_domain_name = "domain2.name"
      t1_renewal_time              = 3600
      t2_rebinding_time            = 3600
    }
  }
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP6_2" {
  name             = "testacc_accessAddAssP6_2"
  routing_instance = junos_routing_instance.testacc_accessAddAssP6.name
  family {
    type    = "inet6"
    network = "fe80:0:0:a::/64"
    dhcp_attributes {
      preferred_lifetime = 7200
      valid_lifetime     = 4800
    }
  }
  link = junos_access_address_assignment_pool.testacc_accessAddAssP6_1.name
}
resource "junos_access_address_assignment_pool" "testacc_accessAddAssP6_3" {
  name             = "testacc_accessAddAssP6_3"
  routing_instance = junos_routing_instance.testacc_accessAddAssP6.name
  family {
    type    = "inet6"
    network = "fe80:0:0:c::/64"
    dhcp_attributes {
      preferred_lifetime_infinite = true
      valid_lifetime_infinite     = true
    }
  }
}
`
}
