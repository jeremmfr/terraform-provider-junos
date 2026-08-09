package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// export TESTACC_INTERFACE=<inteface> to choose interface available else it's ge-0/0/3.
func TestAccResourceRipNeighbor_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
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
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_rip_neighbor.testacc_ripneigh",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_rip_neighbor.testacc_ripneigh2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_rip_neighbor.testacc_ripngneigh",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_rip_neighbor.testacc_ripngneigh2",
					ImportState:       true,
					ImportStateVerify: true,
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

func TestAccResourceRipNeighbor_writeOnly(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					// 1
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckNoResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo",
							"authentication_key"),
						resource.TestCheckNoResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo",
							"authentication_key_wo"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo",
							"authentication_key_wo_version", "1"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.#", "2"),
						resource.TestCheckNoResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.0.key"),
						resource.TestCheckNoResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.0.key_wo"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.0.key_wo_version", "1"),
						resource.TestCheckNoResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.1.key"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.1.key_wo_version", "1"),
					),
				},
				{
					// 2 check that the write-only keys have really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_ripneigh_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`group "?test_rip_group_wo"? neighbor ae0\.0 authentication-key `)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_ripneigh_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`neighbor ae1\.0 authentication-selective-md5 1 key `)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_ripneigh_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`neighbor ae1\.0 authentication-selective-md5 2 key `)),
						),
					},
				},
				{
					// 3
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo",
							"authentication_key_wo_version", "2"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.0.key_wo_version", "2"),
						resource.TestCheckResourceAttr("junos_rip_neighbor.testacc_ripneigh_wo_md5",
							"authentication_selective_md5.1.key_wo_version", "2"),
					),
				},
			},
		})
	}
}
