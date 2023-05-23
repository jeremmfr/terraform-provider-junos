package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosPolicyOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosPolicyOptionsConfigCreate(),
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
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.aggregate_contributor", "true"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.neighbor.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.next_hop.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.prefix_list.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.prefix_list.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.0.route", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.0.option", "exact"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.1.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.1.option", "prefix-length-range"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.1.option_value", "/26-/27"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.neighbor.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.next_hop.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.action", "accept"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.as_path_expand", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.as_path_prepend", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.#", "3"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.0.action", "set"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.0.value", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.1.action", "delete"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.2.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.default_action", "reject"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.load_balance", "per-packet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.local_preference.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.local_preference.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.next", "policy"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.next_hop", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.metric.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.metric.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.preference.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.preference.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.aggregate_contributor", "true"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.neighbor.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.next_hop.*", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.0.route", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.0.option", "exact"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.1.route", "192.0.2.128/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.1.option", "prefix-length-range"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.1.option_value", "/26-/27"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_as_path.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_as_path.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_community.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_community.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.family", "inet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.local_preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.routing_instance", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.interface.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.interface.*", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.neighbor.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.neighbor.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.next_hop.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.next_hop.*", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.policy.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.policy.0", "testacc_policyOptions2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.*", "bgp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.action", "accept"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.as_path_expand", "last-as count 1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.as_path_prepend", "65000 65000"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.#", "3"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.0.action", "set"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.0.value", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.1.action", "delete"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.2.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.default_action", "reject"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.load_balance", "per-packet"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.local_preference.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.local_preference.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.next", "policy"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.next_hop", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.metric.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.metric.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.origin", "igp"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.preference.0.action", "add"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.preference.0.value", "10"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.0.bgp_as_path_group.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.0.bgp_as_path_group.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.0.bgp_as_path_group.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.0.bgp_as_path_group.*", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.local_preference.0.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.metric.0.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.preference.0.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.local_preference.0.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.metric.0.action", "subtract"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.preference.0.action", "subtract"),
					),
				},
				{
					Config: testAccJunosPolicyOptionsConfigUpdate(),
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
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.prefix_list.#", "2"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.*", "ospf"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.metric.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"then.0.preference.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.route_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.#", "1"),
						resource.TestCheckTypeSetElemAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.*", "ospf"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.community.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.local_preference.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.metric.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.then.0.preference.#", "0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.local_preference.0.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.metric.0.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"then.0.preference.0.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.local_preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.local_preference.0.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.metric.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.metric.0.action", "none"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.preference.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"term.0.then.0.preference.0.action", "none"),
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

func testAccJunosPolicyOptionsConfigCreate() string {
	return `
resource "junos_routing_instance" "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource "junos_policyoptions_as_path" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|12|18"
}
resource "junos_policyoptions_as_path_group" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc policyOptions"
    path = "5|12|18"
  }
}
resource "junos_policyoptions_community" "testacc_policyOptions" {
  name    = "testacc_policyOptions"
  members = ["65000:100"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions" {
  name   = "testacc_policyOptions"
  prefix = ["192.0.2.0/25"]
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  from {
    aggregate_contributor = true
    bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_as_path_calc_length {
      count = 4
      match = "orhigher"
    }
    bgp_community = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_community_count {
      count = 6
      match = "orhigher"
    }
    bgp_origin       = "igp"
    color            = 31
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.4"]
    next_hop         = ["192.0.2.4"]
    ospf_area        = "0.0.0.0"
    preference       = 100
    prefix_list      = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
    protocol         = ["bgp"]
    route_filter {
      route  = "192.0.2.0/25"
      option = "exact"
    }
    route_filter {
      route        = "192.0.2.128/25"
      option       = "prefix-length-range"
      option_value = "/26-/27"
    }
  }
  to {
    bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_origin       = "igp"
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.5"]
    next_hop         = ["192.0.2.5"]
    ospf_area        = "0.0.0.0"
    policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
    preference       = 100
    protocol         = ["bgp"]
  }
  then {
    action          = "accept"
    as_path_expand  = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "delete"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "add"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance   = "per-packet"
    local_preference {
      action = "add"
      value  = 10
    }
    next     = "policy"
    next_hop = "192.0.2.4"
    metric {
      action = "add"
      value  = 10
    }
    origin = "igp"
    preference {
      action = "add"
      value  = 10
    }
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_community         = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin            = "igp"
      family                = "inet"
      local_preference      = 100
      routing_instance      = junos_routing_instance.testacc_policyOptions.name
      interface             = ["st0.0"]
      metric                = 5
      neighbor              = ["192.0.2.4"]
      next_hop              = ["192.0.2.4"]
      ospf_area             = "0.0.0.0"
      policy                = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference            = 100
      prefix_list           = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
      protocol              = ["bgp"]
      route_filter {
        route  = "192.0.2.0/25"
        option = "exact"
      }
      route_filter {
        route        = "192.0.2.128/25"
        option       = "prefix-length-range"
        option_value = "/26-/27"
      }
    }
    to {
      bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.5"]
      next_hop         = ["192.0.2.5"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference       = 100
      protocol         = ["bgp"]
    }
    then {
      action          = "accept"
      as_path_expand  = "last-as count 1"
      as_path_prepend = "65000 65000"
      community {
        action = "set"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "delete"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "add"
        value  = junos_policyoptions_community.testacc_policyOptions.name
      }
      default_action = "reject"
      load_balance   = "per-packet"
      local_preference {
        action = "add"
        value  = 10
      }
      next     = "policy"
      next_hop = "192.0.2.4"
      metric {
        action = "add"
        value  = 10
      }
      origin = "igp"
      preference {
        action = "add"
        value  = 10
      }
    }
  }
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions2" {
  name = "testacc_policyOptions2"
  from {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  to {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  then {
    local_preference {
      action = "subtract"
      value  = 10
    }
    metric {
      action = "subtract"
      value  = 10
    }
    preference {
      action = "subtract"
      value  = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "subtract"
        value  = 10
      }
      metric {
        action = "subtract"
        value  = 10
      }
      preference {
        action = "subtract"
        value  = 10
      }
    }
  }
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions3" {
  name                              = "testacc_policyOptions3"
  add_it_to_forwarding_table_export = true
  from {
    route_filter {
      route  = "192.0.2.0/25"
      option = "orlonger"
    }
  }
  then {
    load_balance = "per-packet"
  }
}
`
}

func testAccJunosPolicyOptionsConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource "junos_policyoptions_as_path" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|15"
}
resource "junos_policyoptions_as_path_group" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc_policyOptions"
    path = "5|15"
  }
}
resource "junos_policyoptions_community" "testacc_policyOptions" {
  name    = "testacc_policyOptions"
  members = ["65000:200"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions" {
  name   = "testacc_policyOptions"
  prefix = ["192.0.2.0/26", "192.0.2.64/26"]
}
resource "junos_policyoptions_prefix_list" "testacc_policyOptions2" {
  name       = "testacc_policyOptions2"
  apply_path = "system radius-server <*>"
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions" {
  name = "testacc_policyOptions"
  from {
    aggregate_contributor = true
    bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_as_path_calc_length {
      count = 4
      match = "orhigher"
    }
    bgp_as_path_calc_length {
      count = 3
      match = "equal"
    }
    bgp_as_path_unique_count {
      count = 3
      match = "equal"
    }
    bgp_as_path_unique_count {
      count = 2
      match = "orhigher"
    }
    bgp_community = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_community_count {
      count = 6
      match = "orhigher"
    }
    bgp_community_count {
      count = 5
      match = "equal"
    }
    bgp_origin             = "igp"
    bgp_srte_discriminator = 30

    evpn_esi             = ["00:11:11:11:11:11:11:11:11:33", "00:11:11:11:11:11:11:11:11:32"]
    evpn_mac_route       = "mac-only"
    evpn_tag             = [36, 35, 33]
    family               = "evpn"
    local_preference     = 100
    routing_instance     = junos_routing_instance.testacc_policyOptions.name
    interface            = ["st0.0"]
    metric               = 5
    neighbor             = ["192.0.2.4"]
    next_hop             = ["192.0.2.4"]
    next_hop_type_merged = true
    next_hop_weight {
      match  = "greater-than-equal"
      weight = 500
    }
    next_hop_weight {
      match  = "equal"
      weight = 200
    }
    ospf_area  = "0.0.0.0"
    preference = 100
    prefix_list = [junos_policyoptions_prefix_list.testacc_policyOptions.name,
      junos_policyoptions_prefix_list.testacc_policyOptions2.name,
    ]
    protocol = ["bgp"]
    route_filter {
      route  = "192.0.2.0/25"
      option = "exact"
    }
    route_type          = "internal"
    srte_color          = 39
    state               = "active"
    tunnel_type         = ["ipip"]
    validation_database = "valid"
  }
  to {
    bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
    bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
    bgp_origin       = "igp"
    family           = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface        = ["st0.0"]
    metric           = 5
    neighbor         = ["192.0.2.5"]
    next_hop         = ["192.0.2.5"]
    ospf_area        = "0.0.0.0"
    policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
    preference       = 100
    protocol         = ["ospf"]
  }
  then {
    action          = "accept"
    as_path_expand  = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value  = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance   = "per-packet"
    next           = "policy"
    next_hop       = "192.0.2.4"
    origin         = "igp"
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path           = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_as_path_unique_count {
        count = 4
        match = "orlower"
      }
      bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.4"]
      next_hop         = ["192.0.2.4"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference       = 100
      prefix_list      = [junos_policyoptions_prefix_list.testacc_policyOptions.name]
      protocol         = ["bgp"]
      route_filter {
        route  = "192.0.2.0/25"
        option = "exact"
      }
    }
    to {
      bgp_as_path      = [junos_policyoptions_as_path.testacc_policyOptions.name]
      bgp_community    = [junos_policyoptions_community.testacc_policyOptions.name]
      bgp_origin       = "igp"
      family           = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface        = ["st0.0"]
      metric           = 5
      neighbor         = ["192.0.2.5"]
      next_hop         = ["192.0.2.5"]
      ospf_area        = "0.0.0.0"
      policy           = [junos_policyoptions_policy_statement.testacc_policyOptions2.name]
      preference       = 100
      protocol         = ["ospf"]
    }
    then {
      action          = "accept"
      as_path_expand  = "last-as count 1"
      as_path_prepend = "65000 65000"
      default_action  = "accept"
      load_balance    = "per-packet"
      next            = "policy"
      next_hop        = "192.0.2.4"
      origin          = "igp"
    }
  }
}
resource "junos_policyoptions_policy_statement" "testacc_policyOptions2" {
  name = "testacc_policyOptions2"
  from {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  to {
    bgp_as_path_group = [junos_policyoptions_as_path_group.testacc_policyOptions.name]
  }
  then {
    local_preference {
      action = "none"
      value  = 10
    }
    metric {
      action = "none"
      value  = 10
    }
    preference {
      action = "none"
      value  = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "none"
        value  = 10
      }
      metric {
        action = "none"
        value  = 10
      }
      preference {
        action = "none"
        value  = 10
      }
    }
  }
}
`
}
