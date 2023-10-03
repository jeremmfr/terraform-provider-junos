package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceServicesRpmProbe_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceServicesRpmProbeConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services_rpm_probe.testacc_rpmprobe",
							"test.#", "4"),
					),
				},
				{
					ResourceName:      "junos_services_rpm_probe.testacc_rpmprobe",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceServicesRpmProbeConfigUpdate(),
				},
			},
		})
	}
}

func testAccResourceServicesRpmProbeConfigCreate() string {
	return `
resource "junos_routing_instance" "testacc_rpmprobe" {
  name = "testacc_rpmprobe"
}
resource "junos_services_rpm_probe" "testacc_rpmprobe" {
  name = "testacc_rpmprobe"
  test {
    name                  = "testacc_test#1"
    target_type           = "address"
    target_value          = "192.0.2.33"
    source_address        = "192.0.2.32"
    destination_interface = "ge-0/0/3.0"
  }
  test {
    name                       = "testacc_test#2"
    target_type                = "url"
    target_value               = "https://test.com"
    data_fill                  = "00Aa"
    destination_interface      = "ge-0/0/3.0"
    destination_port           = "7"
    dscp_code_points           = "af11"
    history_size               = 11
    moving_average_size        = 10
    one_way_hardware_timestamp = true
    probe_count                = 15
    probe_interval             = 16
    probe_type                 = "http-get"
    routing_instance           = junos_routing_instance.testacc_rpmprobe.name
    source_address             = "192.0.2.35"
    test_interval              = 23
    thresholds {
      successive_loss = 13
      total_loss      = 14
    }
    traps = ["test-failure", "test-completion"]
    ttl   = 33
  }
  test {
    name                 = "testacc_test#3"
    target_type          = "inet6-address"
    target_value         = "fe80::6"
    probe_type           = "icmp6-ping"
    inet6_source_address = "fe80::1"
  }
  test {
    name               = "testacc_test#4"
    probe_type         = "icmp-ping-timestamp"
    hardware_timestamp = true
    data_size          = 10
    rpm_scale {
      tests_count             = 17
      destination_interface   = "ge-0/0/3.0"
      destination_subunit_cnt = 18
      source_address_base     = "192.0.2.34"
      source_count            = 19
      source_step             = "0.0.0.1"
      target_address_base     = "192.0.2.35"
      target_count            = 21
      target_step             = "0.0.0.1"
    }
    thresholds {
      egress_time     = 24
      ingress_time    = 25
      jitter_egress   = 26
      jitter_ingress  = 27
      jitter_rtt      = 28
      rtt             = 29
      std_dev_egress  = 30
      std_dev_ingress = 31
      std_dev_rtt     = 32
    }
  }
}

resource "junos_services_rpm_probe" "testacc_rpmprobe2" {
  name            = "testacc_rpmprobe2"
  delegate_probes = true
}
`
}

func testAccResourceServicesRpmProbeConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_rpmprobe" {
  name = "testacc_rpmprobe"
}
resource "junos_services_rpm_probe" "testacc_rpmprobe" {
  name = "testacc_rpmprobe"
  test {
    name                  = "testacc_test#4"
    target_type           = "address"
    target_value          = "192.0.2.33"
    source_address        = "192.0.2.32"
    destination_interface = "ge-0/0/3.0"
  }
  test {
    name                       = "testacc_test#3"
    target_type                = "url"
    target_value               = "https://test.com"
    data_fill                  = "00Aa"
    destination_interface      = "ge-0/0/3.0"
    destination_port           = "7"
    dscp_code_points           = "af11"
    history_size               = 11
    moving_average_size        = 10
    one_way_hardware_timestamp = true
    probe_count                = 15
    probe_interval             = 16
    probe_type                 = "http-get"
    routing_instance           = junos_routing_instance.testacc_rpmprobe.name
    source_address             = "192.0.2.35"
    test_interval              = 23
    thresholds {
      successive_loss = 13
      total_loss      = 14
    }
    traps = ["test-failure", "test-completion"]
    ttl   = 33
  }
  test {
    name                 = "testacc_test#2"
    target_type          = "inet6-address"
    target_value         = "fe80::6"
    probe_type           = "icmp6-ping"
    inet6_source_address = "fe80::1"
  }
  test {
    name               = "testacc_test#1"
    probe_type         = "icmp-ping-timestamp"
    hardware_timestamp = true
    data_size          = 10
    rpm_scale {
      tests_count             = 17
      destination_interface   = "ge-0/0/3.0"
      destination_subunit_cnt = 18
      source_address_base     = "192.0.2.34"
      source_count            = 19
      source_step             = "0.0.0.1"
      target_address_base     = "192.0.2.35"
      target_count            = 21
      target_step             = "0.0.0.1"
    }
    thresholds {
      egress_time     = 24
      ingress_time    = 25
      jitter_egress   = 26
      jitter_ingress  = 27
      jitter_rtt      = 28
      rtt             = 29
      std_dev_egress  = 30
      std_dev_ingress = 31
      std_dev_rtt     = 32
    }
  }
}
`
}
