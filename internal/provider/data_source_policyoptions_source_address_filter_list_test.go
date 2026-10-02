package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePolicyoptionsSourceAddressFilterList_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.junos_policyoptions_source_address_filter_list.testacc_dataSAFList",
						"id", "testacc_dataSAFList"),
					resource.TestCheckResourceAttr("data.junos_policyoptions_source_address_filter_list.testacc_dataSAFList",
						"address.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("data.junos_policyoptions_source_address_filter_list.testacc_dataSAFList",
						"address.*", map[string]string{
							"address": "192.0.2.0/25",
							"option":  "orlonger",
						}),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ExpectError:     regexp.MustCompile("policy-options source-address-filter-list .* doesn't exist"),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}
