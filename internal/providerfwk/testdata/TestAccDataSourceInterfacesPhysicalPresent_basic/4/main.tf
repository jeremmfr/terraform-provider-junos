resource "junos_interface_physical" "testacc_dataIfacesPhysPresent" {
  name        = var.interface
  description = "testacc_dataIfacesPhysPresent"
}
resource "junos_interface_logical" "testacc_dataIfacesPhysPresent" {
  name = "${junos_interface_physical.testacc_dataIfacesPhysPresent.name}.0"
  family_inet {}
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentEth" {
  match_name = "^${split("-", var.interface)[0]}-.*$"
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentEth003" {
  match_name = "^${var.interface}$"
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentEth003AdmUp" {
  match_name     = "^${var.interface}$"
  match_admin_up = true
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentEth003OperUp" {
  match_name    = "^${var.interface}$"
  match_oper_up = true
}
data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentLo0" {
  match_name    = "^lo0$"
  match_oper_up = true
}
