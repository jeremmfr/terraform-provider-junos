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

func TestAccResourceSecurityAuthenticationKeyChain_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
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
					ResourceName:      "junos_security_authentication_key_chain.testacc_secauthKeyChain",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_authentication_key_chain.testacc_secauthKeyChainAO",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

// the plan checked as empty by the test framework after each step is what exercises
// keepWriteOnly: without it, the secrets read on the device in Read would be stored
// in the secret arguments of the key blocks and generate a diff.
func TestAccResourceSecurityAuthenticationKeyChain_writeOnly(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_11_0),
			},
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"junos_security_authentication_key_chain.testacc_secauthKeyChain_wo",
							tfjsonpath.New("key_secret_wo"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"5": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"value":   knownvalue.Null(),
									"version": knownvalue.Int64Exact(1),
								}),
								"6": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"value":   knownvalue.Null(),
									"version": knownvalue.Int64Exact(1),
								}),
							}),
						),
						statecheck.ExpectKnownValue(
							"junos_security_authentication_key_chain.testacc_secauthKeyChain_wo",
							tfjsonpath.New("key"),
							knownvalue.SetPartial([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"id":     knownvalue.Int64Exact(4),
									"secret": knownvalue.StringExact("aS3cret#4"),
								}),
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"id":     knownvalue.Int64Exact(5),
									"secret": knownvalue.Null(),
								}),
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"id":     knownvalue.Int64Exact(6),
									"secret": knownvalue.Null(),
								}),
							}),
						),
					},
				},
				{
					// check that the write-only secrets have really been sent to the device
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_secauthKeyChain_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set security authentication-key-chains key-chain "?testacc_secauthKeyChainWO"? key 5 secret `)),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.testacc_secauthKeyChain_wo",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile(
								`set security authentication-key-chains key-chain "?testacc_secauthKeyChainWO"? key 6 secret `)),
						),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"junos_security_authentication_key_chain.testacc_secauthKeyChain_wo",
							tfjsonpath.New("key_secret_wo"),
							knownvalue.MapExact(map[string]knownvalue.Check{
								"5": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"value":   knownvalue.Null(),
									"version": knownvalue.Int64Exact(2),
								}),
								"6": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"value":   knownvalue.Null(),
									"version": knownvalue.Int64Exact(1),
								}),
							}),
						),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"junos_security_authentication_key_chain.testacc_secauthKeyChain_wo",
							tfjsonpath.New("key"),
							knownvalue.SetPartial([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"id":     knownvalue.Int64Exact(4),
									"secret": knownvalue.Null(),
								}),
							}),
						),
					},
				},
			},
		})
	}
}
