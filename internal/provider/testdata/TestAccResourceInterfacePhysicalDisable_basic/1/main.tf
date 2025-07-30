resource "junos_interface_physical" "testacc_interface_disable" {
  name                  = var.interface
  no_disable_on_destroy = true
}
resource "junos_interface_logical" "testacc_interface_disable" {
  name        = "${junos_interface_physical.testacc_interface_disable.name}.0"
  description = "testacc_interface_disable"
}
