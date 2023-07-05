package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceApplications_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceApplicationsPre(),
				},
				{
					Config: testAccDataSourceApplicationsConfig(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_any",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_any",
							"applications.0.name", "any"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh-name",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh-name",
							"applications.0.name", "junos-ssh"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_default_ssh",
							"applications.0.name", "junos-ssh"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_all_ssh",
							"applications.#", "3"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_multi_term",
							"applications.#", "1"),
						resource.TestCheckResourceAttr("data.junos_applications.testacc_multi_term",
							"applications.0.name", "testacc_custom_multi_term"),
					),
				},
			},
		})
	}
}

func testAccDataSourceApplicationsPre() string {
	return `
resource "junos_application" "testacc_custom_ssh" {
  name             = "testacc_custom_ssh"
  protocol         = "tcp"
  destination_port = 22
}
resource "junos_application" "testacc_custom_ssh_term" {
  name = "testacc_custom_ssh_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 22
  }
}
resource "junos_application" "testacc_custom_multi_term" {
  name = "testacc_custom_multi_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 1001
  }
  term {
    name             = "2"
    protocol         = "tcp"
    destination_port = 1002
  }
}
`
}

func testAccDataSourceApplicationsConfig() string {
	return `
resource "junos_application" "testacc_custom_ssh" {
  name             = "testacc_custom_ssh"
  protocol         = "tcp"
  destination_port = 22
}
resource "junos_application" "testacc_custom_ssh_term" {
  name = "testacc_custom_ssh_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 22
  }
}
resource "junos_application" "testacc_custom_multi_term" {
  name = "testacc_custom_multi_term"
  term {
    name             = "1"
    protocol         = "tcp"
    destination_port = 1001
  }
  term {
    name             = "2"
    protocol         = "tcp"
    destination_port = 1002
  }
}

data "junos_applications" "testacc_default_any" {
  match_options {
    protocol = "0"
  }
}
data "junos_applications" "testacc_default_ssh-name" {
  match_name = "^j.*-ssh$"
}
data "junos_applications" "testacc_default_ssh" {
  match_name = "^junos-"
  match_options {
    protocol         = "tcp"
    destination_port = 22
  }
}
data "junos_applications" "testacc_all_ssh" {
  match_options {
    protocol         = "tcp"
    destination_port = 22
  }
}
data "junos_applications" "testacc_multi_term" {
  match_options {
    protocol         = "tcp"
    destination_port = 1001
  }
  match_options {
    protocol         = "tcp"
    destination_port = 1002
  }
}
`
}
