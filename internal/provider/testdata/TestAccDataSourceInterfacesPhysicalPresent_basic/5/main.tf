data "junos_interfaces_physical_present" "testacc_dataIfacesPhysPresentEth003AdmUp" {
  match_name     = "^${var.interface}$"
  match_admin_up = true
}
