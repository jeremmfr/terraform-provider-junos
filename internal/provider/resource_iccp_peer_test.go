package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceIccpPeer_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ResourceName:      "junos_iccp_peer.testacc_iccp_peer",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func TestAccResourceIccpPeer_writeOnly(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key"),
						resource.TestCheckNoResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key_wo"),
						resource.TestCheckResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key_wo_version", "1"),
					),
				},
				{
					// check that the write-only key has really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_iccp_peer_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set protocols iccp peer 192.0.2.2 authentication-key `)),
						),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key"),
						resource.TestCheckNoResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key_wo"),
						resource.TestCheckResourceAttr("junos_iccp_peer.testacc_iccp_peer_wo",
							"authentication_key_wo_version", "2"),
					),
				},
			},
		})
	}
}
