package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosPolicyOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
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
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.0", "192.0.2.0/25"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.aggregate_contributor", "true"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_as_path.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_as_path.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.bgp_community.0", "testacc_policyOptions"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.interface.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.neighbor.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.neighbor.0", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.next_hop.0", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.ospf_area", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.preference", "100"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.prefix_list.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.prefix_list.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.protocol.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.protocol.0", "bgp"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_as_path.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.bgp_community.0", "testacc_policyOptions"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.interface.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.neighbor.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.neighbor.0", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.next_hop.0", "192.0.2.5"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.0", "bgp"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_as_path.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.bgp_community.0", "testacc_policyOptions"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.interface.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.neighbor.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.neighbor.0", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.next_hop.0", "192.0.2.4"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.from.0.protocol.0", "bgp"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_as_path.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_community.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.bgp_community.0", "testacc_policyOptions"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.interface.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.metric", "5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.neighbor.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.neighbor.0", "192.0.2.5"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.next_hop.0", "192.0.2.5"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.0", "bgp"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"from.0.bgp_as_path_group.0", "testacc_policyOptions"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.0.bgp_as_path_group.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions2",
							"to.0.bgp_as_path_group.0", "testacc_policyOptions"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.0", "192.0.2.0/26"),
						resource.TestCheckResourceAttr("junos_policyoptions_prefix_list.testacc_policyOptions",
							"prefix.1", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"from.0.route_filter.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.#", "1"),
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"to.0.protocol.0", "ospf"),
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
						resource.TestCheckResourceAttr("junos_policyoptions_policy_statement.testacc_policyOptions",
							"term.0.to.0.protocol.0", "ospf"),
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
					ResourceName:      "junos_policyoptions_policy_statement.testacc_policyOptions",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosPolicyOptionsConfigCreate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource junos_policyoptions_as_path "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|12|18"
}
resource junos_policyoptions_as_path_group "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc_policyOptions"
    path = "5|12|18"
  }
}
resource junos_policyoptions_community "testacc_policyOptions" {
  name = "testacc_policyOptions"
  members = [ "65000:100" ]
}
resource junos_policyoptions_prefix_list "testacc_policyOptions" {
  name = "testacc_policyOptions"
  prefix = [ "192.0.2.0/25" ]
}
resource junos_policyoptions_policy_statement "testacc_policyOptions" {
  name = "testacc_policyOptions"
  from {
    aggregate_contributor = true
    bgp_as_path = [ junos_policyoptions_as_path.testacc_policyOptions.name ]
    bgp_community = [ junos_policyoptions_community.testacc_policyOptions.name ]
    bgp_origin = "igp"
    family = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface = [ "st0.0" ]
    metric = 5
    neighbor = [ "192.0.2.4" ]
    next_hop = [ "192.0.2.4" ]
    ospf_area = "0.0.0.0"
    preference = 100
    prefix_list = [ junos_policyoptions_prefix_list.testacc_policyOptions.name ]
    protocol = [ "bgp" ]
    route_filter {
      route = "192.0.2.0/25"
      option = "exact"
    }
    route_filter {
      route = "192.0.2.128/25"
      option = "prefix-length-range"
      option_value = "/26-/27"
    }
  }
  to {
    bgp_as_path =  [ junos_policyoptions_as_path.testacc_policyOptions.name ]
    bgp_community = [  junos_policyoptions_community.testacc_policyOptions.name ]
    bgp_origin = "igp"
    family = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface = [ "st0.0" ]
    metric = 5
    neighbor = [ "192.0.2.5" ]
    next_hop = [ "192.0.2.5" ]
    ospf_area = "0.0.0.0"
    policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
    preference = 100
    protocol = [ "bgp" ]
  }
  then {
    action = "accept"
    as_path_expand = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "delete"
      value = junos_policyoptions_community.testacc_policyOptions.name
    }
    community {
      action = "add"
      value = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance = "per-packet"
    local_preference {
      action = "add"
      value = 10
    }
    next = "policy"
    next_hop = "192.0.2.4"
    metric {
      action = "add"
      value = 10
    }
    origin = "igp"
    preference {
      action = "add"
      value = 10
    }
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path = [ junos_policyoptions_as_path.testacc_policyOptions.name ]
      bgp_community = [ junos_policyoptions_community.testacc_policyOptions.name ]
      bgp_origin = "igp"
      family = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface = [ "st0.0" ]
      metric = 5
      neighbor = [ "192.0.2.4" ]
      next_hop = [ "192.0.2.4" ]
      ospf_area = "0.0.0.0"
      policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
      preference = 100
      prefix_list = [ junos_policyoptions_prefix_list.testacc_policyOptions.name ]
      protocol = [ "bgp" ]
      route_filter {
        route = "192.0.2.0/25"
        option = "exact"
      }
      route_filter {
        route = "192.0.2.128/25"
        option = "prefix-length-range"
        option_value = "/26-/27"
      }
    }
    to {
      bgp_as_path =  [ junos_policyoptions_as_path.testacc_policyOptions.name ]
      bgp_community = [  junos_policyoptions_community.testacc_policyOptions.name ]
      bgp_origin = "igp"
      family = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface = [ "st0.0" ]
      metric = 5
      neighbor = [ "192.0.2.5" ]
      next_hop = [ "192.0.2.5" ]
      ospf_area = "0.0.0.0"
      policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
      preference = 100
      protocol = [ "bgp" ]
    }
    then {
      action = "accept"
      as_path_expand = "last-as count 1"
      as_path_prepend = "65000 65000"
      community {
        action = "set"
        value = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "delete"
        value = junos_policyoptions_community.testacc_policyOptions.name
      }
      community {
        action = "add"
        value = junos_policyoptions_community.testacc_policyOptions.name
      }
      default_action = "reject"
      load_balance = "per-packet"
      local_preference {
        action = "add"
        value = 10
      }
      next = "policy"
      next_hop = "192.0.2.4"
      metric {
        action = "add"
        value = 10
      }
      origin = "igp"
      preference {
        action = "add"
        value = 10
      }
    }
  }
}
resource junos_policyoptions_policy_statement "testacc_policyOptions2" {
  name = "testacc_policyOptions2"
  from {
    bgp_as_path_group = [ junos_policyoptions_as_path_group.testacc_policyOptions.name ]
  }
  to {
    bgp_as_path_group = [ junos_policyoptions_as_path_group.testacc_policyOptions.name ]
  }
  then {
    local_preference {
      action = "subtract"
      value = 10
    }
    metric {
      action = "subtract"
      value = 10
    }
    preference {
      action = "subtract"
      value = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "subtract"
        value = 10
      }
      metric {
        action = "subtract"
        value = 10
      }
      preference {
        action = "subtract"
        value = 10
      }
    }
  }
}
`)
}
func testAccJunosPolicyOptionsConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_policyOptions" {
  name = "testacc_policyOptions"
}
resource junos_policyoptions_as_path "testacc_policyOptions" {
  name = "testacc_policyOptions"
  path = "5|15"
}
resource junos_policyoptions_as_path_group "testacc_policyOptions" {
  name = "testacc_policyOptions"
  as_path {
    name = "testacc_policyOptions"
    path = "5|15"
  }
}
resource junos_policyoptions_community "testacc_policyOptions" {
  name = "testacc_policyOptions"
  members = [ "65000:200" ]
}
resource junos_policyoptions_prefix_list "testacc_policyOptions" {
  name = "testacc_policyOptions"
  prefix = [ "192.0.2.0/26", "192.0.2.64/26" ]
}
resource junos_policyoptions_policy_statement "testacc_policyOptions" {
  name = "testacc_policyOptions"
  from {
    aggregate_contributor = true
    bgp_as_path = [ junos_policyoptions_as_path.testacc_policyOptions.name ]
    bgp_community = [ junos_policyoptions_community.testacc_policyOptions.name ]
    bgp_origin = "igp"
    family = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface = [ "st0.0" ]
    metric = 5
    neighbor = [ "192.0.2.4" ]
    next_hop = [ "192.0.2.4" ]
    ospf_area = "0.0.0.0"
    preference = 100
    prefix_list = [ junos_policyoptions_prefix_list.testacc_policyOptions.name ]
    protocol = [ "bgp" ]
    route_filter {
      route = "192.0.2.0/25"
      option = "exact"
    }
  }
  to {
    bgp_as_path =  [ junos_policyoptions_as_path.testacc_policyOptions.name ]
    bgp_community = [  junos_policyoptions_community.testacc_policyOptions.name ]
    bgp_origin = "igp"
    family = "inet"
    local_preference = 100
    routing_instance = junos_routing_instance.testacc_policyOptions.name
    interface = [ "st0.0" ]
    metric = 5
    neighbor = [ "192.0.2.5" ]
    next_hop = [ "192.0.2.5" ]
    ospf_area = "0.0.0.0"
    policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
    preference = 100
    protocol = [ "ospf" ]
  }
  then {
    action = "accept"
    as_path_expand = "65000 65000"
    as_path_prepend = "65000 65000"
    community {
      action = "set"
      value = junos_policyoptions_community.testacc_policyOptions.name
    }
    default_action = "reject"
    load_balance = "per-packet"
    next = "policy"
    next_hop = "192.0.2.4"
    origin = "igp"
  }
  term {
    name = "term"
    from {
      aggregate_contributor = true
      bgp_as_path = [ junos_policyoptions_as_path.testacc_policyOptions.name ]
      bgp_community = [ junos_policyoptions_community.testacc_policyOptions.name ]
      bgp_origin = "igp"
      family = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface = [ "st0.0" ]
      metric = 5
      neighbor = [ "192.0.2.4" ]
      next_hop = [ "192.0.2.4" ]
      ospf_area = "0.0.0.0"
      policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
      preference = 100
      prefix_list = [ junos_policyoptions_prefix_list.testacc_policyOptions.name ]
      protocol = [ "bgp" ]
      route_filter {
        route = "192.0.2.0/25"
        option = "exact"
      }
    }
    to {
      bgp_as_path =  [ junos_policyoptions_as_path.testacc_policyOptions.name ]
      bgp_community = [  junos_policyoptions_community.testacc_policyOptions.name ]
      bgp_origin = "igp"
      family = "inet"
      local_preference = 100
      routing_instance = junos_routing_instance.testacc_policyOptions.name
      interface = [ "st0.0" ]
      metric = 5
      neighbor = [ "192.0.2.5" ]
      next_hop = [ "192.0.2.5" ]
      ospf_area = "0.0.0.0"
      policy = [ junos_policyoptions_policy_statement.testacc_policyOptions2.name ]
      preference = 100
      protocol = [ "ospf" ]
    }
    then {
      action = "accept"
      as_path_expand = "last-as count 1"
      as_path_prepend = "65000 65000"
      default_action = "accept"
      load_balance = "per-packet"
      next = "policy"
      next_hop = "192.0.2.4"
      origin = "igp"
    }
  }
}
resource junos_policyoptions_policy_statement "testacc_policyOptions2" {
  name = "testacc_policyOptions2"
  from {
    bgp_as_path_group = [ junos_policyoptions_as_path_group.testacc_policyOptions.name ]
  }
  to {
    bgp_as_path_group = [ junos_policyoptions_as_path_group.testacc_policyOptions.name ]
  }
  then {
    local_preference {
      action = "none"
      value = 10
    }
    metric {
      action = "none"
      value = 10
    }
    preference {
      action = "none"
      value = 10
    }
    action = "accept"
  }
  term {
    name = "term"
    then {
      local_preference {
        action = "none"
        value = 10
      }
      metric {
        action = "none"
        value = 10
      }
      preference {
        action = "none"
        value = 10
      }
    }
  }
}
`)
}
