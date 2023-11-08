resource "junos_interface_physical" "testacc_dataIfacesPhysPresent" {
  name        = var.interface
  description = "testacc_dataIfacesPhysPresent"
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresent" {
}
