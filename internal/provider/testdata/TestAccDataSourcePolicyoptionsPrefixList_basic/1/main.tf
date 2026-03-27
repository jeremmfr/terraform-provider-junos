resource "junos_policyoptions_prefix_list" "testacc_dataPrefixList" {
  name   = "testacc_dataPrefixList"
  prefix = ["192.0.2.0/25"]
}
