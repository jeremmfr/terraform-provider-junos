resource "junos_interface_physical" "testacc_interface" {
  name        = var.interface
  description = "testacc_interfaceU"
  speed       = "1g"
  ether_opts {
    ae_8023ad = var.interfaceAE
  }
}

