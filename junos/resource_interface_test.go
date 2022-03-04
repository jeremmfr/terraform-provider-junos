package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
// export TESTACC_INTERFACE_AE=ae<num> for choose interface aggregate test else it's ae0.
func TestAccJunosInterface_basic(t *testing.T) {
	if os.Getenv("TESTACC_DEPRECATED") != "" {
		testaccInterface := defaultInterfaceTestAcc
		testaccInterfaceAE := "ae0"
		if os.Getenv("TESTACC_SWITCH") != "" {
			testaccInterface = defaultInterfaceSwitchTestAcc
		}
		if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
			testaccInterface = iface
		}
		if iface := os.Getenv("TESTACC_INTERFACE_AE"); iface != "" {
			testaccInterfaceAE = iface
		}
		if os.Getenv("TESTACC_SWITCH") != "" {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccJunosInterfaceConfigCreate(testaccInterface),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"description", "testacc_interface"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"trunk", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_native", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_members.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_members.0", "100-110"),
						),
					},
					{
						Config: testAccJunosInterfaceConfigUpdate(testaccInterface),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"description", "testacc_interfaceU"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"trunk", "false"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_native", "0"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_members.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"vlan_members.0", "100"),
						),
					},
					{
						ResourceName:      "junos_interface.testacc_interface",
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		} else if os.Getenv("TESTACC_ROUTER") == "" {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccJunosInterfacePlusConfigCreate(testaccInterface, testaccInterfaceAE),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"description", "testacc_interface"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"ether802_3ad", testaccInterfaceAE),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"name", testaccInterfaceAE),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"ae_lacp", "active"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"ae_minimum_links", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"vlan_tagging", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"name", testaccInterfaceAE+".100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"vlan_tagging_id", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"security_zone", "testacc_interface"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"routing_instance", "testacc_interface"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_mtu", "1400"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_filter_input", "testacc_interfaceInet"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_filter_output", "testacc_interfaceInet"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_rpf_check.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.address", "192.0.2.1/25"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.identifier", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.virtual_address.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.virtual_address.0", "192.0.2.2"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.accept_data", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.advertise_interval", "10"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.advertisements_threshold", "3"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.authentication_key", "thePassWord"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.authentication_type", "md5"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.preempt", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.priority", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_interface.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_interface.0.interface", testaccInterfaceAE),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_route.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_mtu", "1400"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_filter_input", "testacc_interfaceInet6"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_filter_output", "testacc_interfaceInet6"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.#", "2"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.1.address", "fe80::1/64"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.address", "2001:db8::1/64"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.identifier", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.virtual_address.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.virtual_address.0", "2001:db8::2"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.virtual_link_local_address", "fe80::2"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.accept_data", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.advertise_interval", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.advertisements_threshold", "3"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.preempt", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.priority", "100"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_interface.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_interface.0.interface", testaccInterfaceAE),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_route.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
						),
					},
					{
						Config: testAccJunosInterfacePlusConfigUpdate(testaccInterface, testaccInterfaceAE),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("junos_interface.testacc_interface",
								"description", "testacc_interfaceU"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"ae_lacp", ""),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAE",
								"ae_minimum_links", "0"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"vlan_tagging_id", "101"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_mtu", "1500"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_mtu", "1500"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_rpf_check.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_rpf_check.0.mode_loose", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.no_accept_data", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.no_preempt", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_interface.#", "0"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet_address.0.vrrp_group.0.track_route.#", "0"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.#", "2"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.#", "1"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.no_accept_data", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.no_preempt", "true"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_interface.#", "0"),
							resource.TestCheckResourceAttr("junos_interface.testacc_interfaceAEunit",
								"inet6_address.0.vrrp_group.0.track_route.#", "0"),
						),
					},
					{
						ResourceName:      "junos_interface.testacc_interface",
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						ResourceName:      "junos_interface.testacc_interfaceAE",
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						ResourceName:      "junos_interface.testacc_interfaceAEunit",
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		}
	}
}

func testAccJunosInterfaceConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interface"
  trunk        = true
  vlan_native  = 100
  vlan_members = ["100-110"]
}
`)
}

func testAccJunosInterfaceConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource junos_interface testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interfaceU"
  vlan_members = ["100"]
}
`)
}

func testAccJunosInterfacePlusConfigCreate(interFace, interfaceAE string) string {
	return fmt.Sprintf(`
resource junos_firewall_filter "testacc_interfaceInet" {
  name   = "testacc_interfaceInet"
  family = "inet"
  term {
    name = "testacc_interface_inetTerm"
    then {
      action = "accept"
    }
  }
}
resource junos_firewall_filter "testacc_interfaceInet6" {
  name   = "testacc_interfaceInet6"
  family = "inet6"
  term {
    name = "testacc_interface_inet6Term"
    then {
      action = "accept"
    }
  }
}
resource junos_security_zone "testacc_interface" {
  name = "testacc_interface"
}
resource junos_routing_instance "testacc_interface" {
  name = "testacc_interface"
}
resource junos_interface testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interface"
  ether802_3ad = "` + interfaceAE + `"
}
resource junos_interface testacc_interfaceAE {
  name             = junos_interface.testacc_interface.ether802_3ad
  description      = "testacc_interfaceAE"
  ae_lacp          = "active"
  ae_minimum_links = 1
  vlan_tagging     = true
}
resource junos_interface testacc_interfaceAEunit {
  name               = "${junos_interface.testacc_interfaceAE.name}.100"
  description        = "testacc_interface_${junos_interface.testacc_interfaceAE.name}.100"
  security_zone      = junos_security_zone.testacc_interface.name
  routing_instance   = junos_routing_instance.testacc_interface.name
  inet_mtu           = 1400
  inet_filter_input  = junos_firewall_filter.testacc_interfaceInet.name
  inet_filter_output = junos_firewall_filter.testacc_interfaceInet.name
  inet_rpf_check {}
  inet_address {
    address = "192.0.2.1/25"
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
        interface     = junos_interface.testacc_interfaceAE.name
        priority_cost = 20
      }
      track_route {
        route            = "192.0.2.128/25"
        routing_instance = "default"
        priority_cost    = 20
      }
    }
  }
  inet6_mtu           = 1400
  inet6_filter_input  = junos_firewall_filter.testacc_interfaceInet6.name
  inet6_filter_output = junos_firewall_filter.testacc_interfaceInet6.name
  inet6_address {
    address = "2001:db8::1/64"
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
        interface     = junos_interface.testacc_interfaceAE.name
        priority_cost = 20
      }
      track_route {
        route            = "192.0.2.128/25"
        routing_instance = "default"
        priority_cost    = 20
      }
    }
  }
  inet6_address {
    address = "fe80::1/64"
  }
}
`)
}

func testAccJunosInterfacePlusConfigUpdate(interFace, interfaceAE string) string {
	return fmt.Sprintf(`
resource junos_firewall_filter "testacc_interfaceInet" {
  name   = "testacc_interfaceInet"
  family = "inet"
  term {
    name = "testacc_interface_inetTerm"
    then {
      action = "accept"
    }
  }
}
resource junos_firewall_filter "testacc_interfaceInet6" {
  name   = "testacc_interfaceInet6"
  family = "inet6"
  term {
    name = "testacc_interface_inet6Term"
    then {
      action = "accept"
    }
  }
}
resource junos_security_zone "testacc_interface" {
  name = "testacc_interface"
}
resource junos_routing_instance "testacc_interface" {
  name = "testacc_interface"
}
resource junos_interface testacc_interface {
  name         = "` + interFace + `"
  description  = "testacc_interfaceU"
  ether802_3ad = "` + interfaceAE + `"
}
resource junos_interface testacc_interfaceAE {
  name         = junos_interface.testacc_interface.ether802_3ad
  description  = "testacc_interfaceAE"
  vlan_tagging = true
}
resource junos_interface testacc_interfaceAEunit {
  name               = "${junos_interface.testacc_interfaceAE.name}.100"
  vlan_tagging_id    = 101
  description        = "testacc_interface_${junos_interface.testacc_interfaceAE.name}.100"
  security_zone      = junos_security_zone.testacc_interface.name
  routing_instance   = junos_routing_instance.testacc_interface.name
  inet_mtu           = 1500
  inet_filter_input  = junos_firewall_filter.testacc_interfaceInet.name
  inet_filter_output = junos_firewall_filter.testacc_interfaceInet.name
  inet_rpf_check {
    mode_loose = true
  }
  inet_address {
    address = "192.0.2.1/25"
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
  inet6_mtu           = 1500
  inet6_filter_input  = junos_firewall_filter.testacc_interfaceInet6.name
  inet6_filter_output = junos_firewall_filter.testacc_interfaceInet6.name
  inet6_address {
    address = "2001:db8::1/64"
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
  inet6_address {
    address = "fe80::1/64"
  }
}
`)
}
