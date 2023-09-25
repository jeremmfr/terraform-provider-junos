package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccResourceInterfaceLogical_basic(t *testing.T) {
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_interface_logical.testacc_interface_logical",
								tfjsonpath.New("family_inet").
									AtMapKey("address").AtSliceIndex(0).
									AtMapKey("vrrp_group").AtSliceIndex(0).
									AtMapKey("authentication_key"),
							),
						},
					},
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
							"family_inet.mtu", "1400"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.filter_input", "testacc_intlogicalInet"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.filter_output", "testacc_intlogicalInet"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.cidr_ip", "192.0.2.1/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.identifier", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.virtual_address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.virtual_address.0", "192.0.2.2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.advertise_interval", "10"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.advertisements_threshold", "3"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.authentication_key", "thePassWord"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.authentication_type", "md5"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.priority", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_interface.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_interface.0.interface", testaccInterface),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_route.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.mtu", "1400"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.filter_input", "testacc_intlogicalInet6"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.filter_output", "testacc_intlogicalInet6"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.#", "2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.1.cidr_ip", "fe80::1/64"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.cidr_ip", "2001:db8::1/64"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.identifier", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.virtual_address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.virtual_address.0", "2001:db8::2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.virtual_link_local_address", "fe80::2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.advertise_interval", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.advertisements_threshold", "3"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.priority", "100"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_interface.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_interface.0.interface", testaccInterface),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_interface.0.priority_cost", "20"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_route.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_route.0.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_route.0.routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_route.0.priority_cost", "20"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"vlan_id", "101"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.mtu", "1500"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.mtu", "1500"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.rpf_check.mode_loose", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.no_accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.no_preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_interface.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet.address.0.vrrp_group.0.track_route.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.#", "2"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.#", "1"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.no_accept_data", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.no_preempt", "true"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_interface.#", "0"),
						resource.TestCheckResourceAttr("junos_interface_logical.testacc_interface_logical",
							"family_inet6.address.0.vrrp_group.0.track_route.#", "0"),
					),
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_interface_logical.testacc_interface_logical",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
			},
		})
	}
}
