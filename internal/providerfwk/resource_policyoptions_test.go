package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyoptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_policyoptions_as_path.testacc_policyOptions",
							"path", "5|12|18"),
						resource.TestCheckResourceAttr("junos_policyoptions_as_path_group.testacc_policyOptions",
							"as_path.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_as_path_group.testacc_policyOptions",
							"as_path.0.path", "5|12|18"),
						resource.TestCheckResourceAttr("junos_policyoptions_community.testacc_policyOptions",
							"members.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_community.testacc_policyOptions",
							"members.0", "65000:100"),
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.*", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.aggregate_contributor", "true"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.neighbor.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.next_hop.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.prefix_list.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.prefix_list.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.0.route", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.0.option", "exact"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.1.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.1.option", "prefix-length-range"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.1.option_value", "/26-/27"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.neighbor.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.next_hop.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.action", "accept"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.as_path_expand", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.as_path_prepend", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.#", "3"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.0.action", "set"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.0.value", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.1.action", "delete"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.2.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.default_action", "reject"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.load_balance", "per-packet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.local_preference.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.local_preference.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.next", "policy"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.next_hop", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.metric.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.metric.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.preference.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.preference.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.aggregate_contributor", "true"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.neighbor.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.next_hop.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.0.route", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.0.option", "exact"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.1.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.1.option", "prefix-length-range"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.1.option_value", "/26-/27"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.neighbor.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.next_hop.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.action", "accept"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.as_path_expand", "last-as count 1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.as_path_prepend", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.#", "3"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.0.action", "set"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.0.value", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.1.action", "delete"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.2.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.default_action", "reject"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.load_balance", "per-packet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.local_preference.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.local_preference.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.next", "policy"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.next_hop", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.metric.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.metric.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.preference.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.preference.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.bgp_as_path_group.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.bgp_as_path_group.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.bgp_as_path_group.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.bgp_as_path_group.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.local_preference.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.metric.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.preference.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.local_preference.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.metric.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.preference.action", "subtract"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_policyoptions_as_path.testacc_policyOptions",
							"path", "5|15"),
						resource.TestCheckResourceAttr("junos_policyoptions_as_path_group.testacc_policyOptions",
							"as_path.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_as_path_group.testacc_policyOptions",
							"as_path.0.path", "5|15"),
						resource.TestCheckResourceAttr("junos_policyoptions_community.testacc_policyOptions",
							"members.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_community.testacc_policyOptions",
							"members.0", "65000:200"),
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.#", "2"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.*", "192.0.2.0/26"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.*", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions2",
							"apply_path", "system radius-server <*>"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.prefix_list.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.route_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.protocol.*", "ospf"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.route_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.protocol.*", "ospf"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.community.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.local_preference.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.metric.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.preference.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.local_preference.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.metric.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.preference.action", "none"),
					),
				},
				{
					ResourceName:      "junos_policyoptions_as_path.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_policyoptions_as_path_group.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_policyoptions_community.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_policyoptions_policy_statement.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_policyoptions_prefix_list.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
