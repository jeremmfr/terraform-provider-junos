resource "junos_interface_physical" "testacc_interface_disable" {
  name                  = var.interface
  description           = "testacc_interface_disable"
  no_disable_on_destroy = true
}
resource "junos_interface_physical_disable" "testacc_interface_disable" {
  name = var.interface
}
