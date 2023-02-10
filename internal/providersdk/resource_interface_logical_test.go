package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosInterfaceLogical_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfaceLogicalConfigCreate(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"description", "testacc_interface_"+testaccInterface+".100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"name", testaccInterface+".100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"vlan_id", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"security_zone", "testacc_interface_logical"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"routing_instance", "testacc_interface_logical"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.mtu", "1400"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.filter_input", "testacc_intlogicalInet"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.filter_output", "testacc_intlogicalInet"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.rpf_check.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.cidr_ip", "192.0.2.1/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.identifier", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.virtual_address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.virtual_address.0", "192.0.2.2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.advertise_interval", "10"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.advertisements_threshold", "3"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.authentication_key", "thePassWord"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.authentication_type", "md5"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.priority", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_interface.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_interface.0.interface", testaccInterface),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_route.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.mtu", "1400"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.filter_input", "testacc_intlogicalInet6"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.filter_output", "testacc_intlogicalInet6"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.#", "2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.1.cidr_ip", "fe80::1/64"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.cidr_ip", "2001:db8::1/64"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.identifier", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.virtual_address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.virtual_address.0", "2001:db8::2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.virtual_link_local_address", "fe80::2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.advertise_interval", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.advertisements_threshold", "3"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.priority", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_interface.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_interface.0.interface", testaccInterface),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_route.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
					),
				},
				{
					Config: testAccJunosInterfaceLogicalConfigUpdate(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"vlan_id", "101"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.mtu", "1500"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.mtu", "1500"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.rpf_check.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.rpf_check.0.mode_loose", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.no_accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.no_preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_interface.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.0.address.0.vrrp_group.0.track_route.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.#", "2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.no_accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.no_preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_interface.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.0.address.0.vrrp_group.0.track_route.#", "0"),
					),
				},
				{
					ResourceName:      "junos_interface_logical.testacc_interface_logical",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosInterfaceLogicalConfigUpdate2(testaccInterface),
				},
				{
					Config: testAccJunosInterfaceLogicalConfigUpdate3(testaccInterface),
				},
				{
					Config: testAccJunosInterfaceLogicalConfigUpdate4(testaccInterface),
				},
			},
		})
	}
}

func testAccJunosInterfaceLogicalConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_firewall_filter" "testacc_intlogicalInet" {
  name   = "testacc_intlogicalInet"
  family = "inet"
  term {
    name = "testacc_intlogicalInetTerm"
    then {
      action = "accept"
    }
  }
}
resource "junos_firewall_filter" "testacc_intlogicalInet6" {
  name   = "testacc_intlogicalInet6"
  family = "inet6"
  term {
    name = "testacc_intlogicalInet6Term"
    then {
      action = "accept"
    }
  }
}
resource "junos_security_zone" "testacc_interface_logical" {
  name = "testacc_interface_logical"
}
resource "junos_routing_instance" "testacc_interface_logical" {
  name = "testacc_interface_logical"
}
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name                       = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  description                = "testacc_interface_${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  disable                    = true
  security_zone              = junos_security_zone.testacc_interface_logical.name
  security_inbound_protocols = ["bgp"]
  security_inbound_services  = ["ssh"]
  routing_instance           = junos_routing_instance.testacc_interface_logical.name
  family_inet {
    mtu           = 1400
    filter_input  = junos_firewall_filter.testacc_intlogicalInet.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet.name
    rpf_check {}
    address {
      cidr_ip = "192.0.2.1/25"
      vrrp_group {
        identifier               = 100
        virtual_address          = ["192.0.2.2"]
        accept_data              = true
        advertise_interval       = 10
        advertisements_threshold = 3
        authentication_key       = "thePassWord"
        authentication_type      = "md5"
        preempt                  = true
        priority                 = 100
        track_interface {
          interface     = junos_interface_physical.testacc_interface_logical_phy.name
          priority_cost = 20
        }
        track_route {
          route            = "192.0.2.128/25"
          routing_instance = "default"
          priority_cost    = 20
        }
      }
    }
  }
  family_inet6 {
    dad_disable   = true
    mtu           = 1400
    filter_input  = junos_firewall_filter.testacc_intlogicalInet6.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet6.name
    address {
      cidr_ip = "2001:db8::1/64"
      vrrp_group {
        identifier                 = 100
        virtual_address            = ["2001:db8::2"]
        virtual_link_local_address = "fe80::2"
        accept_data                = true
        advertise_interval         = 100
        advertisements_threshold   = 3
        preempt                    = true
        priority                   = 100
        track_interface {
          interface     = junos_interface_physical.testacc_interface_logical_phy.name
          priority_cost = 20
        }
        track_route {
          route            = "192.0.2.128/25"
          routing_instance = "default"
          priority_cost    = 20
        }
      }
    }
    address {
      cidr_ip = "fe80::1/64"
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "ip-0/0/0.0"
  tunnel {
    destination         = "192.0.2.10"
    source              = "192.0.2.11"
    allow_fragmentation = true
    path_mtu_discovery  = true
  }
}
`, interFace)
}

func testAccJunosInterfaceLogicalConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_firewall_filter" "testacc_intlogicalInet" {
  name   = "testacc_intlogicalInet"
  family = "inet"
  term {
    name = "testacc_intlogicalInetTerm"
    then {
      action = "accept"
    }
  }
}
resource "junos_firewall_filter" "testacc_intlogicalInet6" {
  name   = "testacc_intlogicalInet6"
  family = "inet6"
  term {
    name = "testacc_intlogicalInet6Term"
    then {
      action = "accept"
    }
  }
}
resource "junos_security_zone" "testacc_interface_logical" {
  name = "testacc_interface"
}
resource "junos_routing_instance" "testacc_interface_logical" {
  name = "testacc_interface"
}
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  lifecycle {
    create_before_destroy = true
  }
  name                       = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id                    = 101
  description                = "testacc_interface_${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  security_zone              = junos_security_zone.testacc_interface_logical.name
  security_inbound_protocols = ["ospf"]
  security_inbound_services  = ["telnet"]
  routing_instance           = junos_routing_instance.testacc_interface_logical.name
  family_inet {
    mtu           = 1500
    filter_input  = junos_firewall_filter.testacc_intlogicalInet.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet.name
    rpf_check {
      mode_loose = true
    }
    address {
      cidr_ip   = "192.0.2.1/25"
      primary   = true
      preferred = true
      vrrp_group {
        identifier               = 100
        virtual_address          = ["192.0.2.2"]
        no_accept_data           = true
        advertise_interval       = 10
        advertisements_threshold = 3
        authentication_key       = "thePassWord"
        authentication_type      = "md5"
        no_preempt               = true
        priority                 = 150
      }
    }
  }
  family_inet6 {
    mtu           = 1500
    filter_input  = junos_firewall_filter.testacc_intlogicalInet6.name
    filter_output = junos_firewall_filter.testacc_intlogicalInet6.name
    address {
      cidr_ip   = "2001:db8::1/64"
      primary   = true
      preferred = true
      vrrp_group {
        identifier                 = 100
        virtual_address            = ["2001:db8::2"]
        virtual_link_local_address = "fe80::2"
        no_accept_data             = true
        advertise_interval         = 100
        no_preempt                 = true
        priority                   = 150
      }
    }
    address {
      cidr_ip = "fe80::1/64"
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "ip-0/0/0.0"
  tunnel {
    destination                  = "192.0.2.12"
    source                       = "192.0.2.13"
    do_not_fragment              = true
    no_path_mtu_discovery        = true
    routing_instance_destination = junos_routing_instance.testacc_interface_logical.name
    traffic_class                = 202
    ttl                          = 203
  }
}
`, interFace)
}

func testAccJunosInterfaceLogicalConfigUpdate2(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name    = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id = 100
  family_inet {
    dhcp {}
  }
  family_inet6 {
    dhcpv6_client {
      client_type                 = "stateful"
      client_identifier_duid_type = "vendor"
      client_ia_type_na           = true
    }
  }
}
`, interFace)
}

func testAccJunosInterfaceLogicalConfigUpdate3(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name    = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id = 100
  family_inet {
    dhcp {
      client_identifier_ascii                        = "BBAA#1"
      client_identifier_prefix_hostname              = true
      client_identifier_prefix_routing_instance_name = true
      client_identifier_use_interface_description    = "device"
      client_identifier_userid_ascii                 = "BBCC#2"
      force_discover                                 = true
      lease_time                                     = 600
      metric                                         = 0
      no_dns_install                                 = true
      options_no_hostname                            = true
      retransmission_attempt                         = 0
      retransmission_interval                        = 4
      server_address                                 = "192.0.2.1"
      update_server                                  = true
      vendor_id                                      = 2
    }
  }
  family_inet6 {
    dhcpv6_client {
      client_type                               = "stateful"
      client_identifier_duid_type               = "vendor"
      client_ia_type_na                         = true
      client_ia_type_pd                         = true
      no_dns_install                            = true
      prefix_delegating_preferred_prefix_length = 0
      prefix_delegating_sub_prefix_length       = 5
      rapid_commit                              = true
      req_option                                = ["fqdn"]
      retransmission_attempt                    = 0
      update_router_advertisement_interface = [
        junos_interface_logical.testacc_interface_logical2.name,
      ]
      update_server = true
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "${junos_interface_physical.testacc_interface_logical_phy.name}.101"
}
`, interFace)
}

func testAccJunosInterfaceLogicalConfigUpdate4(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_interface_logical_phy" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name    = "${junos_interface_physical.testacc_interface_logical_phy.name}.100"
  vlan_id = 100
  family_inet {
    dhcp {
      client_identifier_hexadecimal        = "11BBee"
      client_identifier_userid_hexadecimal = "00AAff"
      lease_time_infinite                  = true
    }
  }
  family_inet6 {
    dhcpv6_client {
      client_type                 = "stateful"
      client_identifier_duid_type = "vendor"
      client_ia_type_pd           = true
      req_option                  = ["fqdn", "domain"]
      update_router_advertisement_interface = [
        junos_interface_logical.testacc_interface_logical2.name,
        "${junos_interface_physical.testacc_interface_logical_phy.name}.100",
      ]
    }
  }
}
resource "junos_interface_logical" "testacc_interface_logical2" {
  name = "${junos_interface_physical.testacc_interface_logical_phy.name}.101"
}
`, interFace)
}
