resource "junos_policyoptions_community" "testacc_dataCommunity" {
  name    = "testacc_dataCommunity"
  members = ["65000:100"]
}
