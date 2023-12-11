resource "junos_interface_logical" "testacc_switchOpts" {
  lifecycle {
    create_before_destroy = true
  }
  name = "lo0.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.16/32"
    }
  }
}
resource "junos_switch_options" "testacc_switchOpts" {
  vtep_source_interface = junos_interface_logical.testacc_switchOpts.name
}
