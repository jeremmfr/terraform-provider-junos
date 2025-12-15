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
)

func TestAccDataSourceConfigRaw_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		// minified not supported
		return
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.junos_config_raw.test_json", "id", "json"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_json", "format", "json"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_json", "config"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_json_minified", "id", "json-minified"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_json_minified", "format", "json-minified"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_json_minified", "config"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_set", "id", "set"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_set", "format", "set"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_set", "config"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_text", "id", "text"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_text", "format", "text"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_text", "config"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_xml", "id", "xml"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_xml", "format", "xml"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_xml", "config"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_xml_minified", "id", "xml-minified"),
					resource.TestCheckResourceAttr("data.junos_config_raw.test_xml_minified", "format", "xml-minified"),
					resource.TestCheckResourceAttrSet("data.junos_config_raw.test_xml_minified", "config"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_json",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("\\s*\"version\"\\s:\\s\".*\",")),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_json_minified",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("\"version\":\".*\",")),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_set",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("set version .*\n")),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_text",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("\\s*version .*;\n")),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_xml",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("<version>.*</version>")),
					),
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.test_xml_minified",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile("<version>.*</version>")),
					),
				},
			},
		},
	})
}

func TestAccDataSourceConfigRaw_switch(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_config_raw.test_json", "id", "json"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_json", "format", "json"),
						resource.TestCheckResourceAttrSet("data.junos_config_raw.test_json", "config"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_set", "id", "set"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_set", "format", "set"),
						resource.TestCheckResourceAttrSet("data.junos_config_raw.test_set", "config"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_text", "id", "text"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_text", "format", "text"),
						resource.TestCheckResourceAttrSet("data.junos_config_raw.test_text", "config"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_xml", "id", "xml"),
						resource.TestCheckResourceAttr("data.junos_config_raw.test_xml", "format", "xml"),
						resource.TestCheckResourceAttrSet("data.junos_config_raw.test_xml", "config"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.test_json",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile("\\s*\"version\"\\s:\\s\".*\",")),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.test_set",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile("set version .*\n")),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.test_text",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile("\\s*version .*;\n")),
						),
						statecheck.ExpectKnownValue(
							"data.junos_config_raw.test_xml",
							tfjsonpath.New("config"),
							knownvalue.StringRegexp(regexp.MustCompile("<version>.*</version>")),
						),
					},
				},
			},
		})
	}
}
