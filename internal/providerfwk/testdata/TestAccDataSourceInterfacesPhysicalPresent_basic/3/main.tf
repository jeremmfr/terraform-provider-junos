resource "junos_interface_physical" "testacc_dataIfacesPhysPresent" {
  name        = var.interface
  description = "testacc_dataIfacesPhysPresent"
}
resource "junos_interface_logical" "testacc_dataIfacesPhysPresent" {
  name = "${junos_interface_physical.testacc_dataIfacesPhysPresent.name}.0"
  family_inet {}
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresent" {
}
